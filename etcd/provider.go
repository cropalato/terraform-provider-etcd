//
// provider.go
// Copyright (C) 2021 rmelo <Ricardo Melo <rmelo@ludia.com>>
//
// Distributed under terms of the MIT license.
//

package etcd

import (
	"context"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/pkg/transport"
)

// Provider -
func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"username": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("ETCD_USERNAME", nil),
			},
			"password": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				DefaultFunc: schema.EnvDefaultFunc("ETCD_PASSWORD", nil),
			},
			"endpoints": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				DefaultFunc: schema.EnvDefaultFunc("ETCD_ENDPOINT", nil),
			},
			"tls": &schema.Schema{
				Type:        schema.TypeBool,
				DefaultFunc: schema.EnvDefaultFunc("ETCD_TLS", true),
				Optional:    true,
				Sensitive:   true,
			},
			"ca_cert": &schema.Schema{
				Type:        schema.TypeString,
				DefaultFunc: schema.EnvDefaultFunc("ETCD_CACERT", nil),
				Optional:    true,
				Sensitive:   true,
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"etcd_key":        resourceKey(),
			"etcd_role":       resourceRole(),
			"etcd_user":       resourceUser(),
			"etcd_permission": resourcePermission(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			"etcd_key":       dataSourceKey(),
			"etcd_keyprefix": dataSourceKeyPrefix(),
		},
		ConfigureContextFunc: providerConfigure,
	}
}

func providerConfigure(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	username := d.Get("username").(string)
	password := d.Get("password").(string)
	tls := d.Get("tls").(bool)
	endpoints := strings.Split(d.Get("endpoints").(string), ",")

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics
	if (username != "") && (password != "") && func(endpointlist []string) bool {
		for _, e := range endpointlist {
			if e == "" {
				return false
			}
		}
		return true
	}(endpoints) == true {
		if tls {
			tlsInfo := transport.TLSInfo{
				TrustedCAFile: d.Get("ca_cert").(string),
			}
			tlsConfig, err := tlsInfo.ClientConfig()
			if err != nil {
				return nil, diag.FromErr(err)
			}
			c, err := clientv3.New(clientv3.Config{
				Endpoints:   endpoints,
				DialTimeout: 5 * time.Second,
				Username:    username,
				Password:    password,
				TLS:         tlsConfig,
			})
			if err != nil {
				return nil, diag.FromErr(err)
			}
			return c, diags
		} else {
			c, err := clientv3.New(clientv3.Config{
				Endpoints:   endpoints,
				DialTimeout: 5 * time.Second,
				Username:    username,
				Password:    password,
			})
			if err != nil {
				return nil, diag.FromErr(err)
			}
			return c, diags
		}

	}

	c, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"localhost:2379"},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		return nil, diag.FromErr(err)
	}

	return c, diags
}
