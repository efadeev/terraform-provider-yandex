package yandex

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/yandex-cloud/go-genproto/yandex/cloud/mdb/mysql/v1"
	"github.com/yandex-cloud/go-genproto/yandex/cloud/operation"
	"google.golang.org/genproto/protobuf/field_mask"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

const (
	yandexMDBMySQLUserCreateTimeout = 10 * time.Minute
	yandexMDBMySQLUserReadTimeout   = 1 * time.Minute
	yandexMDBMySQLUserUpdateTimeout = 10 * time.Minute
	yandexMDBMySQLUserDeleteTimeout = 10 * time.Minute
)

func resourceYandexMDBMySQLUser() *schema.Resource {
	return &schema.Resource{
		Description: "Manages a MySQL user within the Yandex Cloud. For more information, see [the official documentation](https://yandex.cloud/docs/managed-mysql/).",

		Create: resourceYandexMDBMySQLUserCreate,
		Read:   resourceYandexMDBMySQLUserRead,
		Update: resourceYandexMDBMySQLUserUpdate,
		Delete: resourceYandexMDBMySQLUserDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(yandexMDBMySQLUserCreateTimeout),
			Read:   schema.DefaultTimeout(yandexMDBMySQLUserReadTimeout),
			Update: schema.DefaultTimeout(yandexMDBMySQLUserUpdateTimeout),
			Delete: schema.DefaultTimeout(yandexMDBMySQLUserDeleteTimeout),
		},

		SchemaVersion: 0,

		Schema: map[string]*schema.Schema{
			"cluster_id": {
				Type:        schema.TypeString,
				Description: "The ID of the MySQL cluster.",
				Required:    true,
				ForceNew:    true,
			},
			"name": {
				Type:        schema.TypeString,
				Description: "The name of the user.",
				Required:    true,
			},
			"password": {
				Type:        schema.TypeString,
				Description: "The password of the user.",
				Optional:    true,
				Sensitive:   true,
			},
			"permission": {
				Type:        schema.TypeSet,
				Description: "Set of permissions granted to the user.",
				Optional:    true,
				Computed:    true,
				Set:         mysqlUserPermissionHash,
				Elem:        resourceYandexMDBMySQLUserPermission(),
			},
			"global_permissions": {
				Type:        schema.TypeSet,
				Description: "List user's global permissions. Allowed permissions: `REPLICATION_CLIENT`, `REPLICATION_SLAVE`, `PROCESS`, `FLUSH_OPTIMIZER_COSTS`, `SHOW_ROUTINE`, `MDB_ADMIN` for clear list use empty list. If the attribute is not specified there will be no changes.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Optional: true,
				Computed: true,
			},
			"connection_limits": {
				Type:        schema.TypeList,
				Description: "User's connection limits. If the attribute is not specified there will be no changes. Default value is `-1`. When these parameters are set to `-1`, backend default values will be actually used.",
				Optional:    true,
				Computed:    true,
				MaxItems:    1,
				Elem:        resourceYandexMDBMySQLUserConnectionLimits(),
			},
			"authentication_plugin": {
				Type:        schema.TypeString,
				Description: "Authentication plugin. Allowed values: `MYSQL_NATIVE_PASSWORD`, `CACHING_SHA2_PASSWORD`, `SHA256_PASSWORD` (for version 5.7 `MYSQL_NATIVE_PASSWORD`, `SHA256_PASSWORD`)",
				Optional:    true,
				Computed:    true,
			},
			"connection_manager": {
				Type:        schema.TypeMap,
				Description: "Connection Manager connection configuration. Filled in by the server automatically.",
				Computed:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"generate_password": {
				Type:        schema.TypeBool,
				Description: "Generate password using Connection Manager. Allowed values: `true` or `false`. It's used only during user creation and is ignored during updating.\n\n~> **Must specify either password or generate_password**.\n",
				Optional:    true,
				Default:     false,
			},
		},
	}
}

func resourceYandexMDBMySQLUserPermission() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"database_name": {
				Type:        schema.TypeString,
				Description: "The name of the database that the permission grants access to.",
				Required:    true,
			},
			"roles": {
				Type:        schema.TypeList,
				Description: "List user's roles in the database. Allowed roles: `ALL`,`ALTER`,`ALTER_ROUTINE`,`CREATE`,`CREATE_ROUTINE`,`CREATE_TEMPORARY_TABLES`, `CREATE_VIEW`,`DELETE`,`DROP`,`EVENT`,`EXECUTE`,`INDEX`,`INSERT`,`LOCK_TABLES`,`SELECT`,`SHOW_VIEW`,`TRIGGER`,`UPDATE`.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Optional: true,
			},
		},
	}
}

func resourceYandexMDBMySQLUserConnectionLimits() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"max_questions_per_hour": {
				Type:        schema.TypeInt,
				Description: "Max questions per hour.",
				Optional:    true,
				Default:     -1,
			},
			"max_updates_per_hour": {
				Type:        schema.TypeInt,
				Description: "Max updates per hour.",
				Optional:    true,
				Default:     -1,
			},
			"max_connections_per_hour": {
				Type:        schema.TypeInt,
				Description: "Max connections per hour.",
				Optional:    true,
				Default:     -1,
			},
			"max_user_connections": {
				Type:        schema.TypeInt,
				Description: "Max user connections.",
				Optional:    true,
				Default:     -1,
			},
		},
	}
}

func resourceYandexMDBMySQLUserCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	ctx, cancel := config.ContextWithTimeout(d.Timeout(schema.TimeoutUpdate))
	defer cancel()

	clusterID := d.Get("cluster_id").(string)
	userSpec, err := expandMySQLUserSpec(d)
	if err != nil {
		return err
	}

	if !isValidMySQLPasswordConfiguration(userSpec) {
		return fmt.Errorf("must specify either password or generate_password")
	}

	request := &mysql.CreateUserRequest{
		ClusterId: clusterID,
		UserSpec:  userSpec,
	}
	op, err := retryConflictingOperation(ctx, config, func() (*operation.Operation, error) {
		log.Printf("[DEBUG] Sending MySQL user create request: %+v", request)
		return config.sdk.MDB().MySQL().User().Create(ctx, request)
	})

	userID := constructResourceId(clusterID, userSpec.Name)
	d.SetId(userID)

	if err != nil {
		return fmt.Errorf("error while requesting API to create user for MySQL Cluster %q: %s", clusterID, err)
	}

	if err = op.Wait(ctx); err != nil {
		return fmt.Errorf("error while creating user for MySQL Cluster %q: %s", clusterID, err)
	}

	if _, err := op.Response(); err != nil {
		return fmt.Errorf("creating user for MySQL Cluster %q failed: %s", clusterID, err)
	}

	return resourceYandexMDBMySQLUserRead(d, meta)
}

func expandMySQLUserSpec(d *schema.ResourceData) (*mysql.UserSpec, error) {
	user := &mysql.UserSpec{}

	if v, ok := d.GetOk("name"); ok {
		user.Name = v.(string)
	}

	if v, ok := d.GetOk("password"); ok {
		user.Password = v.(string)
	}

	if v, ok := d.GetOk("permission"); ok {
		permissions, err := expandMysqlUserPermissions(v.(*schema.Set))
		if err != nil {
			return nil, err
		}
		user.Permissions = permissions
	}

	if v, ok := d.GetOk("global_permissions"); ok {
		gs, err := expandMysqlUserGlobalPermissions(v.(*schema.Set).List())
		if err != nil {
			return nil, err
		}
		user.GlobalPermissions = gs
	}

	if conLimits, ok := d.GetOk("connection_limits"); ok {
		connectionLimitsMap := (conLimits.([]interface{}))[0].(map[string]interface{})
		user.ConnectionLimits = expandMySQLConnectionLimits(connectionLimitsMap)
	}

	if v, ok := d.GetOk("authentication_plugin"); ok {
		authenticationPlugin, err := expandEnum("authentication_plugin", v.(string), mysql.AuthPlugin_value)
		if err != nil {
			return nil, err
		}
		user.AuthenticationPlugin = mysql.AuthPlugin(*authenticationPlugin)
	}

	if v, ok := d.GetOk("generate_password"); ok {
		user.GeneratePassword = wrapperspb.Bool(v.(bool))
	}

	return user, nil
}

func resourceYandexMDBMySQLUserRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	ctx, cancel := config.ContextWithTimeout(d.Timeout(schema.TimeoutRead))
	defer cancel()

	clusterID, username, err := deconstructResourceId(d.Id())
	if err != nil {
		return err
	}

	user, err := config.sdk.MDB().MySQL().User().Get(ctx, &mysql.GetUserRequest{
		ClusterId: clusterID,
		UserName:  username,
	})
	if err != nil {
		return handleNotFoundError(err, d, fmt.Sprintf("User %q", username))
	}

	permissions, err := flattenMysqlUserPermissions(user.Permissions)
	if err != nil {
		return err
	}

	connectionLimits := flattenMysqlUserConnectionLimits(user)
	globalPermissions := unbindGlobalPermissions(user.GlobalPermissions)

	d.Set("cluster_id", clusterID)
	d.Set("name", user.Name)
	d.Set("permission", permissions)
	d.Set("global_permissions", globalPermissions)
	d.Set("connection_limits", connectionLimits)
	if user.AuthenticationPlugin != 0 {
		d.Set("authentication_plugin", mysql.AuthPlugin_name[int32(user.AuthenticationPlugin)])
	}
	d.Set("connection_manager", flattenMySQLUserConnectionManager(user.ConnectionManager))
	return nil
}

func resourceYandexMDBMySQLUserUpdate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	ctx, cancel := config.ContextWithTimeout(d.Timeout(schema.TimeoutDelete))
	defer cancel()

	user, err := expandMySQLUserSpec(d)
	if err != nil {
		return err
	}

	if !isValidMySQLPasswordConfiguration(user) {
		return fmt.Errorf("must specify either password or generate_password")
	}

	clusterID := d.Get("cluster_id").(string)
	request := &mysql.UpdateUserRequest{
		ClusterId:            clusterID,
		UserName:             user.Name,
		Password:             user.Password,
		Permissions:          user.Permissions,
		AuthenticationPlugin: user.AuthenticationPlugin,
		ConnectionLimits:     user.ConnectionLimits,
		GlobalPermissions:    user.GlobalPermissions,
		UpdateMask:           &field_mask.FieldMask{Paths: []string{"authentication_plugin", "password", "permissions", "connection_limits", "global_permissions"}},
	}
	op, err := retryConflictingOperation(ctx, config, func() (*operation.Operation, error) {
		log.Printf("[DEBUG] Sending MySQL user update request: %+v", request)
		return config.sdk.MDB().MySQL().User().Update(ctx, request)
	})
	if err != nil {
		return fmt.Errorf("error while requesting API to update user in MySQL Cluster %q: %s", clusterID, err)
	}

	if err = op.Wait(ctx); err != nil {
		return fmt.Errorf("error while updating user in MySQL Cluster %q: %s", clusterID, err)
	}

	if _, err := op.Response(); err != nil {
		return fmt.Errorf("updating user for MySQL Cluster %q failed: %s", clusterID, err)
	}
	return nil
}

func resourceYandexMDBMySQLUserDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	ctx, cancel := config.ContextWithTimeout(d.Timeout(schema.TimeoutDelete))
	defer cancel()

	clusterID := d.Get("cluster_id").(string)
	username := d.Get("name").(string)

	request := &mysql.DeleteUserRequest{
		ClusterId: clusterID,
		UserName:  username,
	}
	op, err := retryConflictingOperation(ctx, config, func() (*operation.Operation, error) {
		log.Printf("[DEBUG] Sending MySQL user delete request: %+v", request)
		return config.sdk.MDB().MySQL().User().Delete(ctx, request)
	})
	if err != nil {
		return fmt.Errorf("error while requesting API to delete user from MySQL Cluster %q: %s", clusterID, err)
	}

	if err = op.Wait(ctx); err != nil {
		return fmt.Errorf("error while deleting user from MySQL Cluster %q: %s", clusterID, err)
	}

	if _, err := op.Response(); err != nil {
		return fmt.Errorf("deleting user from MySQL Cluster %q failed: %s", clusterID, err)
	}

	return nil
}
