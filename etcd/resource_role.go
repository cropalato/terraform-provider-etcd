package etcd

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"go.etcd.io/etcd/client/v3"
	//"go.etcd.io/etcd/client/v3/concurrency"
)

func resourceRole() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceRoleCreate,
		ReadContext:   resourceRoleRead,
		UpdateContext: resourceRoleUpdate,
		DeleteContext: resourceRoleDelete,
		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			//"key": &schema.Schema{
			//	Type:     schema.TypeString,
			//	Optional: true,
			//},
			//"withPrefix": &schema.Schema{
			//	Type:     schema.TypeBool,
			//	Optional: true,
			//},
			//"permission": &schema.Schema{
			//	Type:     schema.TypeInt,
			//	Optional: true,
			//},
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func resourceRoleCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	var diags diag.Diagnostics
	var requestTimeout = 5 * time.Second

	cli := m.(*clientv3.Client)

	name := fmt.Sprintf("%v", d.Get("name"))
	//key := fmt.Sprintf("%v", d.Get("key"))
	//withPrefix := fmt.Sprintf("%v", d.Get("withPrefix"))
	//permission := fmt.Sprintf("%v", d.Get("permission"))
	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	_, err := cli.RoleGet(ctx, name)
	cancel()
	if err == nil {
		return append(diag.FromErr(err), diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "ResourceRoleCreate.",
			Detail:   fmt.Sprintf("The role %v already exist and it is not managed by this terraform.", name),
		})
	}
	ctx, cancel = context.WithTimeout(context.Background(), requestTimeout)
	_, err = cli.RoleAdd(ctx, name)
	cancel()
	if err != nil {
		return diag.FromErr(err)
	}
	// always run
	d.SetId(name)

	resourceRoleRead(ctx, d, m)

	return diags
}

func resourceRoleRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	var diags diag.Diagnostics
	//var requestTimeout = 5 * time.Second

	//cli := m.(*clientv3.Client)

	//session, err := concurrency.NewSession(cli)
	//if err != nil {
	//	return diag.FromErr(err)
	//}
	//defer session.Close()

	//ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)

	//key := fmt.Sprintf("%v", d.Get("key"))
	//if key == "" {
	//	key = string(d.Id())
	//	d.Set("key", key)
	//}
	//m1 := concurrency.NewMutex(session, fmt.Sprintf("/resourceKeyRead/%v", key))
	//if err := m1.Lock(context.TODO()); err != nil {
	//	return diag.FromErr(err)
	//}
	//resp, err := cli.Get(ctx, key)
	//if err := m1.Unlock(context.TODO()); err != nil {
	//	return diag.FromErr(err)
	//}
	//cancel()
	//if err != nil {
	//	return append(diag.FromErr(err), diag.Diagnostic{
	//		Severity: diag.Error,
	//		Summary:  "Error reading data from etcd",
	//		Detail:   fmt.Sprintf("Failed calling cli.Get() from resourceKeyRead(). key=/%v/", key),
	//	})
	//}
	//if resp.Count == 0 {
	//	return append(diag.FromErr(err), diag.Diagnostic{
	//		Severity: diag.Error,
	//		Summary:  "Error reading data from etcd",
	//		Detail:   "Etcd returns no answer. It is suppose to have at least one empty value.",
	//	})
	//}
	//for _, ev := range resp.Kvs {
	//	if err := d.Set("value", string(ev.Value)); err != nil {
	//		return append(diag.FromErr(err), diag.Diagnostic{
	//			Severity: diag.Error,
	//			Summary:  "Error reading data from etcd server",
	//			Detail:   "Failed saving data into 'value'.",
	//		})
	//	}
	//	//if err := d.Set("create_revision", int(ev.CreateRevision)); err != nil {
	//	//	return append(diag.FromErr(err), diag.Diagnostic{
	//	//		Severity: diag.Error,
	//	//		Summary:  "Error reading data from etcd server",
	//	//		Detail:   "Failed saving data into 'create_revision'.",
	//	//	})
	//	//}
	//	//if err := d.Set("mod_revision", int(ev.ModRevision)); err != nil {
	//	//	return append(diag.FromErr(err), diag.Diagnostic{
	//	//		Severity: diag.Error,
	//	//		Summary:  "Error reading data from etcd server",
	//	//		Detail:   "Failed saving data into 'mod_revision'.",
	//	//	})
	//	//}
	//	//if err := d.Set("version", int(ev.Version)); err != nil {
	//	//	return append(diag.FromErr(err), diag.Diagnostic{
	//	//		Severity: diag.Error,
	//	//		Summary:  "Error reading data from etcd server",
	//	//		Detail:   "Failed saving data into 'version'.",
	//	//	})
	//	//}
	//	break
	//}

	//// always run
	//d.SetId(key)

	return diags
}
func resourceRoleUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	var diags diag.Diagnostics
	//var requestTimeout = 5 * time.Second
	//if d.HasChange("value") {

	//	cli := m.(*clientv3.Client)

	//	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	//	key := fmt.Sprintf("%v", d.Get("key"))
	//	value := fmt.Sprintf("%v", d.Get("value"))

	//	session, err := concurrency.NewSession(cli)
	//	if err != nil {
	//		return diag.FromErr(err)
	//	}
	//	defer session.Close()

	//	m1 := concurrency.NewMutex(session, fmt.Sprintf("/resourceKeyRead/%v", key))
	//	if err := m1.Lock(context.TODO()); err != nil {
	//		return diag.FromErr(err)
	//	}

	//	if err := m1.Lock(context.TODO()); err != nil {
	//		return diag.FromErr(err)
	//	}
	//	// Should I do a Get()?
	//	resp, err := cli.Get(ctx, key)
	//	if err := m1.Unlock(context.TODO()); err != nil {
	//		return diag.FromErr(err)
	//	}

	//	cancel()
	//	if err != nil {
	//		return append(diag.FromErr(err), diag.Diagnostic{
	//			Severity: diag.Error,
	//			Summary:  "Error Creating resource Key.",
	//			Detail:   "Error Calling Get() funcion from resourceKeyCreate.",
	//		})
	//	}
	//	if resp.Count == 0 {
	//		return append(diag.FromErr(err), diag.Diagnostic{
	//			Severity: diag.Error,
	//			Summary:  "Filed Creating resource Key.",
	//			Detail:   "The Key already exists",
	//		})
	//	}

	//	ctx, cancel = context.WithTimeout(context.Background(), requestTimeout)
	//	if err := m1.Lock(context.TODO()); err != nil {
	//		return diag.FromErr(err)
	//	}
	//	_, put_err := cli.Put(ctx, key, value)
	//	if err := m1.Unlock(context.TODO()); err != nil {
	//		return diag.FromErr(err)
	//	}
	//	cancel()
	//	if put_err != nil {
	//		return append(diag.FromErr(err), diag.Diagnostic{
	//			Severity: diag.Error,
	//			Summary:  "Failed Creating resource Key.",
	//			Detail:   "Error writing key/value in etcd server",
	//		})
	//	}

	//	resourceKeyRead(ctx, d, m)

	//	//d.Set("last_updated", time.Now().Format(time.RFC850))
	//}

	resourceRoleRead(ctx, d, m)
	return diags

}

func resourceRoleDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	var diags diag.Diagnostics
	//var requestTimeout = 5 * time.Second
	//cli := m.(*clientv3.Client)

	//ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	//key := fmt.Sprintf("%v", d.Get("key"))

	//session, err := concurrency.NewSession(cli)
	//if err != nil {
	//	return diag.FromErr(err)
	//}
	//defer session.Close()

	//m1 := concurrency.NewMutex(session, fmt.Sprintf("/resourceKeyRead/%v", key))
	//if err := m1.Lock(context.TODO()); err != nil {
	//	return diag.FromErr(err)
	//}

	//if err := m1.Lock(context.TODO()); err != nil {
	//	return diag.FromErr(err)
	//}
	//// Should I do a Get()?
	//_, derr := cli.Delete(ctx, key)
	//if err := m1.Unlock(context.TODO()); err != nil {
	//	return diag.FromErr(err)
	//}
	//cancel()
	//if derr != nil {
	//	return append(diag.FromErr(err), diag.Diagnostic{
	//		Severity: diag.Error,
	//		Summary:  "Error Deleting resource Key.",
	//		Detail:   "Error Calling Delete() funcion from resourceKeyCreate.",
	//	})
	//}

	return diags
}
