package etcd

import (
	"context"
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"go.etcd.io/etcd/client/v3"
	//"go.etcd.io/etcd/client/v3/concurrency"
)

var (
	lowerCharSet   = "abcdefghijklmnopqrstuvwxyz"
	upperCharSet   = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	specialCharSet = "!@#$%&*+-_?.,"
	numberSet      = "0123456789"
	allCharSet     = lowerCharSet + upperCharSet + specialCharSet + numberSet
)

func resourceUser() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceUserCreate,
		ReadContext:   resourceUserRead,
		UpdateContext: resourceUserUpdate,
		DeleteContext: resourceUserDelete,
		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"password": &schema.Schema{
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

// this function has been copied from https://golangbyexample.com/generate-random-password-golang/
func generatePassword(passwdLen, nSpecialChar, nNum, nUpperCase int) string {
	var password strings.Builder

	//Set special character
	for i := 0; i < nSpecialChar; i++ {
		random := rand.Intn(len(specialCharSet))
		password.WriteString(string(specialCharSet[random]))
	}

	//Set numeric
	for i := 0; i < nNum; i++ {
		random := rand.Intn(len(numberSet))
		password.WriteString(string(numberSet[random]))
	}

	//Set uppercase
	for i := 0; i < nUpperCase; i++ {
		random := rand.Intn(len(upperCharSet))
		password.WriteString(string(upperCharSet[random]))
	}

	remainingLength := passwdLen - nSpecialChar - nNum - nUpperCase
	for i := 0; i < remainingLength; i++ {
		random := rand.Intn(len(allCharSet))
		password.WriteString(string(allCharSet[random]))
	}
	inRune := []rune(password.String())
	rand.Shuffle(len(inRune), func(i, j int) {
		inRune[i], inRune[j] = inRune[j], inRune[i]
	})
	return string(inRune)
}

func resourceUserCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	var diags diag.Diagnostics
	var requestTimeout = 5 * time.Second

	cli := m.(*clientv3.Client)

	name := fmt.Sprintf("%v", d.Get("name"))
	password := fmt.Sprintf("%v", d.Get("password"))
	if password == "" {
		password = generatePassword(24, 3, 3, 3)
		d.Set("password", password)
	}
	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	_, err := cli.UserAdd(ctx, name, password)
	cancel()
	if err != nil {
		return diag.FromErr(err)
	}
	// always run
	d.SetId(uuidGenerator())

	resourceUserRead(ctx, d, m)

	return diags
}

func resourceUserRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	var diags diag.Diagnostics
	var requestTimeout = 5 * time.Second

	cli := m.(*clientv3.Client)

	name := fmt.Sprintf("%v", d.Get("name"))
	if name == "" {
		name = string(d.Id())
		d.Set("name", name)
	}
	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	_, err := cli.UserGet(ctx, name)
	cancel()
	if err != nil {
		return append(diag.FromErr(err), diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "ResourceUserRead.",
			Detail:   fmt.Sprintf("The user %v doesn't exist. Maybe someone removed that manually.", name),
		})
	}
	// always run
	d.SetId(uuidGenerator())

	return diags
}
func resourceUserUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	var diags diag.Diagnostics
	//var requestTimeout = 5 * time.Second

	//old_value, new_value := d.GetChange("name")

	//cli := m.(*clientv3.Client)

	//ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	//role, err := cli.RoleGet(ctx, fmt.Sprintf("%v", old_value))
	//cancel()
	//if err != nil {
	//	return append(diag.FromErr(err), diag.Diagnostic{
	//		Severity: diag.Error,
	//		Summary:  "ResourceRoleUpdate.",
	//		Detail:   fmt.Sprintf("The original role %v doesn't exist.", old_value),
	//	})
	//}
	//ctx, cancel = context.WithTimeout(context.Background(), requestTimeout)
	//_, err = cli.RoleGet(ctx, fmt.Sprintf("%v", new_value))
	//cancel()
	//if err == nil {
	//	return append(diag.FromErr(err), diag.Diagnostic{
	//		Severity: diag.Error,
	//		Summary:  "ResourceRoleUpdate.",
	//		Detail:   fmt.Sprintf("The new role name %v already exist.", new_value),
	//	})
	//}
	//// Creating new role
	//ctx, cancel = context.WithTimeout(context.Background(), requestTimeout)
	//_, err = cli.RoleAdd(ctx, fmt.Sprintf("%v", new_value))
	//cancel()
	//if err != nil {
	//	return diag.FromErr(err)
	//}
	//for _, p := range role.Perm {
	//	ctx, cancel = context.WithTimeout(context.Background(), requestTimeout)
	//	if fmt.Sprintf("%v", p.PermType) == "READWRITE" {
	//		_, err = cli.RoleGrantPermission(ctx, fmt.Sprintf("%v", new_value), string(p.Key), string(p.RangeEnd), clientv3.PermissionType(clientv3.PermReadWrite))
	//	} else {
	//		_, err = cli.RoleGrantPermission(ctx, fmt.Sprintf("%v", new_value), string(p.Key), string(p.RangeEnd), clientv3.PermissionType(clientv3.PermRead))
	//	}
	//	cancel()
	//	if err != nil {
	//		return append(diags, diag.Diagnostic{
	//			Severity: diag.Error,
	//			Summary:  "ResourceRoleUpdate.",
	//			Detail:   fmt.Sprintf("Failed copying grants from old role %v to new role %v.", old_value, new_value),
	//		})
	//	}
	//}
	//ctx, cancel = context.WithTimeout(context.Background(), requestTimeout)
	//_, err = cli.RoleDelete(ctx, fmt.Sprintf("%v", old_value))
	//cancel()
	//if err != nil {
	//	return diag.FromErr(err)
	//}

	//// always run
	//d.Set("name", fmt.Sprintf("%v", new_value))
	//d.SetId(fmt.Sprintf("%v", new_value))
	resourceUserRead(ctx, d, m)
	return diags

}

func resourceUserDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	var diags diag.Diagnostics
	var requestTimeout = 5 * time.Second

	cli := m.(*clientv3.Client)

	name := fmt.Sprintf("%v", d.Get("name"))
	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	_, err := cli.UserDelete(ctx, name)
	cancel()
	if err != nil {
		return diag.FromErr(err)
	}

	return diags
}
