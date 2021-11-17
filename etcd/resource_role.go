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
	d.SetId(uuidGenerator())

	resourceRoleRead(ctx, d, m)

	return diags
}

func resourceRoleRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	var diags diag.Diagnostics
	var requestTimeout = 5 * time.Second

	cli := m.(*clientv3.Client)

	name := fmt.Sprintf("%v", d.Get("name"))
	if name == "" {
		name = string(d.Id())
		d.Set("name", name)
	}
	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	_, err := cli.RoleGet(ctx, name)
	cancel()
	if err != nil {
		return append(diag.FromErr(err), diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "ResourceRoleRead.",
			Detail:   fmt.Sprintf("The role %v doesn't exist. Maybe someone removed that manually.", name),
		})
	}
	// always run
	d.SetId(uuidGenerator())

	return diags
}
func resourceRoleUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	var diags diag.Diagnostics
	var requestTimeout = 5 * time.Second

	old_value, new_value := d.GetChange("name")

	cli := m.(*clientv3.Client)

	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	role, err := cli.RoleGet(ctx, fmt.Sprintf("%v", old_value))
	cancel()
	if err != nil {
		return append(diag.FromErr(err), diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "ResourceRoleUpdate.",
			Detail:   fmt.Sprintf("The original role %v doesn't exist.", old_value),
		})
	}
	ctx, cancel = context.WithTimeout(context.Background(), requestTimeout)
	_, err = cli.RoleGet(ctx, fmt.Sprintf("%v", new_value))
	cancel()
	if err == nil {
		return append(diag.FromErr(err), diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "ResourceRoleUpdate.",
			Detail:   fmt.Sprintf("The new role name %v already exist.", new_value),
		})
	}
	// Creating new role
	ctx, cancel = context.WithTimeout(context.Background(), requestTimeout)
	_, err = cli.RoleAdd(ctx, fmt.Sprintf("%v", new_value))
	cancel()
	if err != nil {
		return diag.FromErr(err)
	}
	for _, p := range role.Perm {
		ctx, cancel = context.WithTimeout(context.Background(), requestTimeout)
		if fmt.Sprintf("%v", p.PermType) == "READWRITE" {
			_, err = cli.RoleGrantPermission(ctx, fmt.Sprintf("%v", new_value), string(p.Key), string(p.RangeEnd), clientv3.PermissionType(clientv3.PermReadWrite))
		} else {
			_, err = cli.RoleGrantPermission(ctx, fmt.Sprintf("%v", new_value), string(p.Key), string(p.RangeEnd), clientv3.PermissionType(clientv3.PermRead))
		}
		cancel()
		if err != nil {
			return append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "ResourceRoleUpdate.",
				Detail:   fmt.Sprintf("Failed copying grants from old role %v to new role %v.", old_value, new_value),
			})
		}
	}
	ctx, cancel = context.WithTimeout(context.Background(), requestTimeout)
	_, err = cli.RoleDelete(ctx, fmt.Sprintf("%v", old_value))
	cancel()
	if err != nil {
		return diag.FromErr(err)
	}

	// always run
	d.Set("name", fmt.Sprintf("%v", new_value))
	d.SetId(uuidGenerator())
	resourceRoleRead(ctx, d, m)
	return diags

}

func resourceRoleDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

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
	if err != nil {
		return append(diag.FromErr(err), diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "ResourceRoleDelete.",
			Detail:   fmt.Sprintf("The role %v doesn't exist.", name),
		})
	}
	ctx, cancel = context.WithTimeout(context.Background(), requestTimeout)
	_, err = cli.RoleDelete(ctx, name)
	cancel()
	if err != nil {
		return diag.FromErr(err)
	}

	return diags
}
