package etcd

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"go.etcd.io/etcd/client/v3"
)

type Etcd_key struct {
	Key             string
	Value           string
	Create_revision int64
	Mod_revision    int64
	Version         int64
}

func dataSourceKey() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceKeyRead,
		Schema: map[string]*schema.Schema{
			"endpoints": &schema.Schema{
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"value": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"create_revision": &schema.Schema{
							Type:     schema.TypeInt,
							Computed: true,
						},
						"mod_revision": &schema.Schema{
							Type:     schema.TypeInt,
							Computed: true,
						},
						"version": &schema.Schema{
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceKeyRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	c := m.(*clientv3.Client)
	if c == c {
		fmt.Println("ops")
	}

	var diags diag.Diagnostics
	var requestTimeout = 5 * time.Second

	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"localhost:2379"},
		DialTimeout: requestTimeout,
		Username:    "rmelo",
		Password:    "rmelo",
	})
	if err != nil {
		return append(diag.FromErr(err), diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to create HashiCups client2",
			Detail:   "Unable to auth user for authenticated HashiCups client3",
		})
		//return diag.FromErr(err)
	}
	defer cli.Close()

	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	resp, err := cli.Get(ctx, "/foo")
	cancel()
	if err != nil {
		return append(diag.FromErr(err), diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to create HashiCups client2",
			Detail:   "Unable to auth user for authenticated HashiCups client3",
		})
	}
	if resp.Count == 0 {
		d.Set("value", "WTF")
		return diags
	}
	objmaparray := make([]map[string]interface{}, 0)
	for _, ev := range resp.Kvs {
		key := fmt.Sprintf("[{\"key\": \"%s\", \"value\": \"%s\", \"create_revision\": %d, \"mod_revision\": %d, \"version\": %d}]", string(ev.Key),
			string(ev.Value), int(ev.CreateRevision),
			int(ev.ModRevision), int(ev.Version))
		err := json.NewDecoder(strings.NewReader(key)).Decode(&objmaparray)
		if err != nil {
			return diag.FromErr(err)
		}
		if err := d.Set("endpoints", objmaparray); err != nil {
			return append(diag.FromErr(err), diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Unable to create HashiCups client2",
				Detail:   "Unable to auth user for authenticated HashiCups client3",
			})
		}
		break
	}

	// always run
	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))

	return diags
}
