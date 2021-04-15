package etcd

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"go.etcd.io/etcd/client/v3"
)

type Etcd_key struct {
	Key   string
	Value string
	//Create_revision int64
	//Mod_revision    int64
	//Version         int64
}

func dataSourceKey() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceKeyRead,
		Schema: map[string]*schema.Schema{
			"key": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"value": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			//"create_revision": &schema.Schema{
			//	Type:     schema.TypeInt,
			//	Computed: true,
			//},
			//"mod_revision": &schema.Schema{
			//	Type:     schema.TypeInt,
			//	Computed: true,
			//},
			//"version": &schema.Schema{
			//	Type:     schema.TypeInt,
			//	Computed: true,
			//},
			"last_updated": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceKeyRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	var diags diag.Diagnostics
	var requestTimeout = 5 * time.Second

	cli := m.(*clientv3.Client)

	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	key := fmt.Sprintf("%v", d.Get("key"))
	resp, err := cli.Get(ctx, key)
	cancel()
	if err != nil {
		return append(diag.FromErr(err), diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error reading data from etcd",
			Detail:   "Failed calling cli.Get() from dataSourceKeyRead()",
		})
	}
	if resp.Count == 0 {
		return append(diag.FromErr(err), diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error reading data from etcd",
			Detail:   "Etcd returns no answer. It is suppose to have at least one empty value.",
		})
	}
	for _, ev := range resp.Kvs {
		if err := d.Set("value", string(ev.Value)); err != nil {
			return append(diag.FromErr(err), diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Error reading data from etcd server",
				Detail:   "Failed saving data into 'value'.",
			})
		}
		//if err := d.Set("create_revision", int(ev.CreateRevision)); err != nil {
		//	return append(diag.FromErr(err), diag.Diagnostic{
		//		Severity: diag.Error,
		//		Summary:  "Error reading data from etcd server",
		//		Detail:   "Failed saving data into 'create_revision'.",
		//	})
		//}
		//if err := d.Set("mod_revision", int(ev.ModRevision)); err != nil {
		//	return append(diag.FromErr(err), diag.Diagnostic{
		//		Severity: diag.Error,
		//		Summary:  "Error reading data from etcd server",
		//		Detail:   "Failed saving data into 'mod_revision'.",
		//	})
		//}
		//if err := d.Set("version", int(ev.Version)); err != nil {
		//	return append(diag.FromErr(err), diag.Diagnostic{
		//		Severity: diag.Error,
		//		Summary:  "Error reading data from etcd server",
		//		Detail:   "Failed saving data into 'version'.",
		//	})
		//}
		break
	}

	// always run
	d.SetId(key)

	return diags
}
