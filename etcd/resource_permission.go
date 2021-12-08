package etcd

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	clientv3 "go.etcd.io/etcd/client/v3"
)

func resourcePermission() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourcePermissionCreate,
		ReadContext:   resourcePermissionRead,
		UpdateContext: resourcePermissionUpdate,
		DeleteContext: resourcePermissionDelete,
		Schema: map[string]*schema.Schema{
			"role": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"key": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"withprefix": &schema.Schema{
				Type:     schema.TypeBool,
				Required: true,
			},
			"endrange": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},
			"permission": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
					v := val.(string)
					if !contains([]string{"READ", "READWRITE"}, v) {
						errs = append(errs, fmt.Errorf("%q must be READ or READWRITE, got: %v", key, v))
					}
					return
				},
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func resourcePermissionCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	var diags diag.Diagnostics
	var requestTimeout = 5 * time.Second
	var rangeEnd string
	var permission clientv3.PermissionType

	cli := m.(*clientv3.Client)
	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)

	role := fmt.Sprintf("%v", d.Get("role"))
	key := fmt.Sprintf("%v", d.Get("key"))
	withPrefix := d.Get("withprefix").(bool)
	if withPrefix == true {
		rangeEnd = clientv3.GetPrefixRangeEnd(key)
	} else {
		rangeEnd = d.Get("endrange").(string)
		if rangeEnd == "" {
			return append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "resourcePermissionCreate error.",
				Detail:   fmt.Sprintf("'endrange' is a mandatory argument when you define 'withprefix' == false."),
			})
		}
	}
	if fmt.Sprintf("%v", d.Get("permission")) == "READWRITE" {
		permission = clientv3.PermissionType(clientv3.PermReadWrite)
	} else {
		permission = clientv3.PermissionType(clientv3.PermRead)
	}
	_, err := cli.RoleGrantPermission(ctx, role, key, rangeEnd, permission)
	cancel()
	if err != nil {
		return append(diag.FromErr(err), diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "resourcePermissionCreate error.",
			Detail:   fmt.Sprintf("Failed creating permission: %v to key: %v into role: %v", permission, key, role),
		})
	}

	// always run
	d.SetId(uuidGenerator())

	resourceUserRead(ctx, d, m)

	return diags
}

func resourcePermissionRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	var diags diag.Diagnostics
	var requestTimeout = 5 * time.Second

	cli := m.(*clientv3.Client)
	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)

	role := fmt.Sprintf("%v", d.Get("role"))
	key := fmt.Sprintf("%v", d.Get("key"))
	resp, err := cli.RoleGet(ctx, role)
	cancel()
	if err != nil {
		return append(diag.FromErr(err), diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "resourcePermissionRead error.",
			Detail:   fmt.Sprintf("Failed getting role: %v", role),
		})
	}
	for _, p := range resp.Perm {
		if string(p.Key) != key {
			continue
		}
		if string(p.RangeEnd) == clientv3.GetPrefixRangeEnd(key) {
			d.Set("withPrefix", true)
		} else {
			d.Set("withPrefix", false)
		}
		d.Set("permission", fmt.Sprintf("%v", p.PermType))
		d.SetId(uuidGenerator())
	}

	return diags
}

func resourcePermissionUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	var diags diag.Diagnostics
	var requestTimeout = 5 * time.Second
	var rangeEnd string
	var permission clientv3.PermissionType

	cli := m.(*clientv3.Client)
	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)

	role := fmt.Sprintf("%v", d.Get("role"))
	key := fmt.Sprintf("%v", d.Get("key"))
	withPrefix := d.Get("withprefix").(bool)
	if withPrefix == true {
		rangeEnd = clientv3.GetPrefixRangeEnd(key)
	} else {
		rangeEnd = d.Get("endrange").(string)
		if rangeEnd == "" {
			return append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "resourcePermissionUpdate error.",
				Detail:   fmt.Sprintf("'endrange' is a mandatory argument when you define 'withprefix' == false."),
			})
		}
	}
	if fmt.Sprintf("%v", d.Get("permission")) == "READWRITE" {
		permission = clientv3.PermissionType(clientv3.PermReadWrite)
	} else {
		permission = clientv3.PermissionType(clientv3.PermRead)
	}
	_, err := cli.RoleGrantPermission(ctx, role, key, rangeEnd, permission)
	cancel()
	if err != nil {
		return append(diag.FromErr(err), diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Resource Permission error.",
			Detail:   fmt.Sprintf("Failed creating permission: %v to key: %v into role: %v", permission, key, role),
		})
	}

	// always run
	d.SetId(uuidGenerator())

	return diags

}

func resourcePermissionDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	var diags diag.Diagnostics
	var requestTimeout = 5 * time.Second

	cli := m.(*clientv3.Client)
	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)

	role := fmt.Sprintf("%v", d.Get("role"))
	key := fmt.Sprintf("%v", d.Get("key"))
	rangeEnd := d.Get("endrange").(string)

	resp, err := cli.RoleGet(ctx, role)
	cancel()
	if err != nil {
		return append(diag.FromErr(err), diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "resourcePermissionDelete error.",
			Detail:   fmt.Sprintf("Failed getting role: %v", role),
		})
	}
	for _, p := range resp.Perm {
		if string(p.Key) != key || string(p.RangeEnd) != rangeEnd {
			continue
		}
		ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
		_, err = cli.RoleRevokePermission(ctx, role, key, rangeEnd)
		cancel()
		if err != nil {
			return append(diag.FromErr(err), diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "resourcePermissionDelete error.",
				Detail:   fmt.Sprintf("Failed revoking permission to key: %v from role: %v", key, role),
			})
		}
	}

	return diags
}
