package etcd

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/concurrency"
)

func resourceKey() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceKeyCreate,
		ReadContext:   resourceKeyRead,
		UpdateContext: resourceKeyUpdate,
		DeleteContext: resourceKeyDelete,
		Schema: map[string]*schema.Schema{
			"key": &schema.Schema{
				Description: "Etcd key",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"value": &schema.Schema{
				Description: "Etcd value",
				Type:        schema.TypeString,
				Optional:    true,
			},
			//"create_revision": &schema.Schema{
			//	Type:     schema.TypeInt,
			//	Computed: true,
			//	Optional: true,
			//},
			//"mod_revision": &schema.Schema{
			//	Type:     schema.TypeInt,
			//	Computed: true,
			//	Optional: true,
			//},
			//"version": &schema.Schema{
			//	Type:     schema.TypeInt,
			//	Computed: true,
			//	Optional: true,
			//},
			//"last_updated": &schema.Schema{
			//	Type:     schema.TypeString,
			//	Optional: true,
			//	Computed: true,
			//},
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func resourceKeyCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	var diags diag.Diagnostics
	var requestTimeout = 5 * time.Second

	cli := m.(*clientv3.Client)

	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	key := fmt.Sprintf("%v", d.Get("key"))
	value := fmt.Sprintf("%v", d.Get("value"))
	resp, err := cli.Get(ctx, key)
	cancel()
	if err != nil {
		return append(diag.FromErr(err), diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error Creating resource Key.",
			Detail:   "Error Calling Get() funcion from resourceKeyCreate.",
		})
	}
	if resp.Count == 0 {
		return append(diag.FromErr(err), diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Filed Creating resource Key.",
			Detail:   "The Key already exists",
		})
	}

	ctx, cancel = context.WithTimeout(context.Background(), requestTimeout)
	_, put_err := cli.Put(ctx, key, value)
	cancel()
	if put_err != nil {
		return append(diag.FromErr(err), diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Failed Creating resource Key.",
			Detail:   "Error writing key/value in etcd server",
		})
	}
	// always run
	d.SetId(uuidGenerator())

	resourceKeyRead(ctx, d, m)

	return diags
}

func resourceKeyRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	var diags diag.Diagnostics
	var requestTimeout = 5 * time.Second

	cli := m.(*clientv3.Client)

	session, err := concurrency.NewSession(cli)
	if err != nil {
		return diag.FromErr(err)
	}
	defer session.Close()

	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)

	key := fmt.Sprintf("%v", d.Get("key"))
	if key == "" {
		key = string(d.Id())
		d.Set("key", key)
	}
	m1 := concurrency.NewMutex(session, fmt.Sprintf("/resourceKeyRead/%v", key))
	if err := m1.Lock(context.TODO()); err != nil {
		return diag.FromErr(err)
	}
	resp, err := cli.Get(ctx, key)
	if err := m1.Unlock(context.TODO()); err != nil {
		return diag.FromErr(err)
	}
	cancel()
	if err != nil {
		return append(diag.FromErr(err), diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error reading data from etcd",
			Detail:   fmt.Sprintf("Failed calling cli.Get() from resourceKeyRead(). key=/%v/", key),
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
	d.SetId(uuidGenerator())

	return diags
}
func resourceKeyUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	var diags diag.Diagnostics
	var requestTimeout = 5 * time.Second
	if d.HasChange("value") {

		cli := m.(*clientv3.Client)

		ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
		key := fmt.Sprintf("%v", d.Get("key"))
		value := fmt.Sprintf("%v", d.Get("value"))

		session, err := concurrency.NewSession(cli)
		if err != nil {
			return diag.FromErr(err)
		}
		defer session.Close()

		m1 := concurrency.NewMutex(session, fmt.Sprintf("/resourceKeyRead/%v", key))
		if err := m1.Lock(context.TODO()); err != nil {
			return diag.FromErr(err)
		}

		if err := m1.Lock(context.TODO()); err != nil {
			return diag.FromErr(err)
		}
		// Should I do a Get()?
		resp, err := cli.Get(ctx, key)
		if err := m1.Unlock(context.TODO()); err != nil {
			return diag.FromErr(err)
		}

		cancel()
		if err != nil {
			return append(diag.FromErr(err), diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Error Creating resource Key.",
				Detail:   "Error Calling Get() funcion from resourceKeyCreate.",
			})
		}
		if resp.Count == 0 {
			return append(diag.FromErr(err), diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Filed Creating resource Key.",
				Detail:   "The Key already exists",
			})
		}

		ctx, cancel = context.WithTimeout(context.Background(), requestTimeout)
		if err := m1.Lock(context.TODO()); err != nil {
			return diag.FromErr(err)
		}
		_, put_err := cli.Put(ctx, key, value)
		if err := m1.Unlock(context.TODO()); err != nil {
			return diag.FromErr(err)
		}
		cancel()
		if put_err != nil {
			return append(diag.FromErr(err), diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Failed Creating resource Key.",
				Detail:   "Error writing key/value in etcd server",
			})
		}

		resourceKeyRead(ctx, d, m)

		//d.Set("last_updated", time.Now().Format(time.RFC850))
	}

	return diags

}

func resourceKeyDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	var diags diag.Diagnostics
	var requestTimeout = 5 * time.Second
	cli := m.(*clientv3.Client)

	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	key := fmt.Sprintf("%v", d.Get("key"))

	session, err := concurrency.NewSession(cli)
	if err != nil {
		return diag.FromErr(err)
	}
	defer session.Close()

	m1 := concurrency.NewMutex(session, fmt.Sprintf("/resourceKeyRead/%v", key))
	if err := m1.Lock(context.TODO()); err != nil {
		return diag.FromErr(err)
	}

	if err := m1.Lock(context.TODO()); err != nil {
		return diag.FromErr(err)
	}
	// Should I do a Get()?
	_, derr := cli.Delete(ctx, key)
	if err := m1.Unlock(context.TODO()); err != nil {
		return diag.FromErr(err)
	}
	cancel()
	if derr != nil {
		return append(diag.FromErr(err), diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error Deleting resource Key.",
			Detail:   "Error Calling Delete() funcion from resourceKeyCreate.",
		})
	}

	return diags
}
