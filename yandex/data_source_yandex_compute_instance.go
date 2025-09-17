package yandex

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/yandex-cloud/go-genproto/yandex/cloud/compute/v1"
	"github.com/yandex-cloud/go-sdk/sdkresolvers"
	"github.com/yandex-cloud/terraform-provider-yandex/common"
)

func dataSourceYandexComputeInstance() *schema.Resource {
	return &schema.Resource{
		Description: "Get information about a Yandex Compute instance. For more information, see [the official documentation](https://yandex.cloud/docs/compute/concepts/vm).\n\n~> One of `instance_id` or `name` should be specified.\n",

		Read: dataSourceYandexComputeInstanceRead,
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:        schema.TypeString,
				Description: "The ID of a specific instance.",
				Optional:    true,
				Computed:    true,
			},
			"name": {
				Type:        schema.TypeString,
				Description: common.ResourceDescriptions["name"],
				Optional:    true,
				Computed:    true,
			},
			"fqdn": {
				Type:        schema.TypeString,
				Description: resourceYandexComputeInstance().Schema["fqdn"].Description,
				Computed:    true,
			},
			"folder_id": {
				Type:        schema.TypeString,
				Description: common.ResourceDescriptions["folder_id"],
				Optional:    true,
				Computed:    true,
			},
			"zone": {
				Type:        schema.TypeString,
				Description: common.ResourceDescriptions["zone"],
				Computed:    true,
			},
			"description": {
				Type:        schema.TypeString,
				Description: common.ResourceDescriptions["description"],
				Computed:    true,
			},
			"labels": {
				Type:        schema.TypeMap,
				Description: common.ResourceDescriptions["labels"],
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Set:         schema.HashString,
			},
			"metadata": {
				Type:        schema.TypeMap,
				Description: resourceYandexComputeInstance().Schema["metadata"].Description,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Set:         schema.HashString,
			},
			"platform_id": {
				Type:        schema.TypeString,
				Description: resourceYandexComputeInstance().Schema["platform_id"].Description,
				Computed:    true,
			},
			"status": {
				Type:        schema.TypeString,
				Description: resourceYandexComputeInstance().Schema["status"].Description,
				Computed:    true,
			},
			"resources": {
				Type:        schema.TypeList,
				Description: resourceYandexComputeInstance().Schema["resources"].Description,
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"memory": {
							Type:     schema.TypeFloat,
							Computed: true,
						},
						"cores": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"gpus": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"core_fraction": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},
			"boot_disk": {
				Type:        schema.TypeList,
				Description: resourceYandexComputeInstance().Schema["boot_disk"].Description,
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"auto_delete": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"device_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"mode": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"disk_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"initialize_params": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"description": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"size": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"block_size": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"type": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"image_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"snapshot_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"kms_key_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
					},
				},
			},
			"network_acceleration_type": {
				Type:        schema.TypeString,
				Description: resourceYandexComputeInstance().Schema["network_acceleration_type"].Description,
				Computed:    true,
			},
			"network_interface": {
				Type:        schema.TypeList,
				Description: resourceYandexComputeInstance().Schema["network_interface"].Description,
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"index": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"mac_address": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"ipv4": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"ip_address": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"ipv6": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"ipv6_address": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"subnet_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"nat": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"nat_ip_address": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"nat_ip_version": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"security_group_ids": {
							Type:     schema.TypeSet,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
							Set:      schema.HashString,
						},
						"dns_record": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"fqdn": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"dns_zone_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"ttl": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"ptr": {
										Type:     schema.TypeBool,
										Computed: true,
									},
								},
							},
						},
						"ipv6_dns_record": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"fqdn": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"dns_zone_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"ttl": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"ptr": {
										Type:     schema.TypeBool,
										Computed: true,
									},
								},
							},
						},
						"nat_dns_record": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"fqdn": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"dns_zone_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"ttl": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"ptr": {
										Type:     schema.TypeBool,
										Computed: true,
									},
								},
							},
						},
					},
				},
			},
			"secondary_disk": {
				Type:        schema.TypeSet,
				Description: resourceYandexComputeInstance().Schema["secondary_disk"].Description,
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"auto_delete": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"device_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"mode": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"disk_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"scheduling_policy": {
				Type:        schema.TypeList,
				Description: resourceYandexComputeInstance().Schema["scheduling_policy"].Description,
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"preemptible": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  false,
						},
					},
				},
			},
			"service_account_id": {
				Type:        schema.TypeString,
				Description: common.ResourceDescriptions["service_account_id"],
				Optional:    true,
			},

			"placement_policy": {
				Type:        schema.TypeList,
				Description: resourceYandexComputeInstance().Schema["placement_policy"].Description,
				MaxItems:    1,
				Optional:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"placement_group_id": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"placement_group_partition": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"host_affinity_rules": {
							Type:       schema.TypeList,
							Computed:   true,
							Optional:   true,
							ConfigMode: schema.SchemaConfigModeAttr,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"key": {
										Type:     schema.TypeString,
										Required: true,
									},
									"op": {
										Type:     schema.TypeString,
										Required: true,
									},
									"values": {
										Type:     schema.TypeList,
										Required: true,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
								},
							},
						},
					},
				},
			},

			"created_at": {
				Type:        schema.TypeString,
				Description: common.ResourceDescriptions["created_at"],
				Computed:    true,
			},

			"local_disk": {
				Type:        schema.TypeList,
				Description: resourceYandexComputeInstance().Schema["local_disk"].Description,
				Optional:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"size_bytes": {
							Type:     schema.TypeInt,
							Required: true,
							ForceNew: true,
						},
						"device_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},

			"metadata_options": {
				Type:        schema.TypeList,
				Description: resourceYandexComputeInstance().Schema["metadata_options"].Description,
				MaxItems:    1,
				Optional:    true,
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"gce_http_endpoint": {
							Type:         schema.TypeInt,
							ValidateFunc: validation.IntBetween(0, 2),
							Optional:     true,
							Computed:     true,
						},
						"aws_v1_http_endpoint": {
							Type:         schema.TypeInt,
							ValidateFunc: validation.IntBetween(0, 2),
							Optional:     true,
							Computed:     true,
						},
						"gce_http_token": {
							Type:         schema.TypeInt,
							ValidateFunc: validation.IntBetween(0, 2),
							Optional:     true,
							Computed:     true,
						},
						"aws_v1_http_token": {
							Type:         schema.TypeInt,
							ValidateFunc: validation.IntBetween(0, 2),
							Optional:     true,
							Computed:     true,
						},
					},
				},
			},

			"filesystem": {
				Type:        schema.TypeSet,
				Description: resourceYandexComputeInstance().Schema["filesystem"].Description,
				Optional:    true,
				Set:         hashFilesystem,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"filesystem_id": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"device_name": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"mode": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
					},
				},
			},

			"gpu_cluster_id": {
				Type:        schema.TypeString,
				Description: resourceYandexComputeInstance().Schema["gpu_cluster_id"].Description,
				Optional:    true,
				Computed:    true,
			},

			"maintenance_policy": {
				Type:        schema.TypeString,
				Description: resourceYandexComputeInstance().Schema["maintenance_policy"].Description,
				Optional:    true,
				Computed:    true,
			},
			"maintenance_grace_period": {
				Type:        schema.TypeString,
				Description: resourceYandexComputeInstance().Schema["maintenance_grace_period"].Description,
				Optional:    true,
				Computed:    true,
			},

			"hardware_generation": {
				Type:        schema.TypeList,
				Description: resourceYandexComputeInstance().Schema["hardware_generation"].Description,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"legacy_features": {
							Type: schema.TypeList,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"pci_topology": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
							Computed: true,
						},

						"generation2_features": {
							Type: schema.TypeList,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{},
							},
							Computed: true,
						},
					},
				},
				Computed: true,
			},
		},
	}

}

func dataSourceYandexComputeInstanceRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	ctx := config.Context()

	err := checkOneOf(d, "instance_id", "name")
	if err != nil {
		return err
	}

	instanceID := d.Get("instance_id").(string)
	_, instanceNameOk := d.GetOk("name")

	if instanceNameOk {
		instanceID, err = resolveObjectID(ctx, config, d, sdkresolvers.InstanceResolver)
		if err != nil {
			return fmt.Errorf("failed to resolve data source instance by name: %v", err)
		}
	}

	instance, err := config.sdk.Compute().Instance().Get(ctx, &compute.GetInstanceRequest{
		InstanceId: instanceID,
		View:       compute.InstanceView_FULL,
	})

	if err != nil {
		return handleNotFoundError(err, d, fmt.Sprintf("instance with ID %q", instanceID))
	}

	resources, err := flattenInstanceResources(instance)
	if err != nil {
		return err
	}

	bootDisk, err := flattenInstanceBootDisk(ctx, instance, config.sdk.Compute().Disk())
	if err != nil {
		return err
	}

	networkInterfaces, _, _, err := flattenInstanceNetworkInterfaces(instance)
	if err != nil {
		return err
	}

	secondaryDisks, err := flattenInstanceSecondaryDisks(instance)
	if err != nil {
		return err
	}

	schedulingPolicy, err := flattenInstanceSchedulingPolicy(instance)
	if err != nil {
		return err
	}

	placementPolicy, err := flattenInstancePlacementPolicy(instance)
	if err != nil {
		return err
	}

	localDisks := flattenLocalDisks(instance)

	metadataOptions := flattenInstanceMetadataOptions(instance)

	filesystems := flattenInstanceFilesystems(instance)

	hardwareGeneration, err := flattenComputeHardwareGeneration(instance.HardwareGeneration)
	if err != nil {
		return err
	}

	d.Set("created_at", getTimestamp(instance.CreatedAt))
	d.Set("instance_id", instance.Id)
	d.Set("platform_id", instance.PlatformId)
	d.Set("folder_id", instance.FolderId)
	d.Set("zone", instance.ZoneId)
	d.Set("name", instance.Name)
	d.Set("fqdn", instance.Fqdn)
	d.Set("description", instance.Description)
	d.Set("service_account_id", instance.ServiceAccountId)
	d.Set("status", strings.ToLower(instance.Status.String()))
	d.Set("metadata_options", metadataOptions)

	if err := d.Set("metadata", instance.Metadata); err != nil {
		return err
	}

	if err := d.Set("labels", instance.Labels); err != nil {
		return err
	}

	if err := d.Set("resources", resources); err != nil {
		return err
	}

	if err := d.Set("boot_disk", bootDisk); err != nil {
		return err
	}

	if instance.NetworkSettings != nil {
		d.Set("network_acceleration_type", strings.ToLower(instance.NetworkSettings.Type.String()))
	}

	if err := d.Set("network_interface", networkInterfaces); err != nil {
		return err
	}

	if err := d.Set("secondary_disk", secondaryDisks); err != nil {
		return err
	}

	if err := d.Set("scheduling_policy", schedulingPolicy); err != nil {
		return err
	}

	if err := d.Set("placement_policy", placementPolicy); err != nil {
		return err
	}

	if err := d.Set("local_disk", localDisks); err != nil {
		return err
	}

	if err := d.Set("filesystem", filesystems); err != nil {
		return err
	}

	if instance.GpuSettings != nil {
		d.Set("gpu_cluster_id", instance.GpuSettings.GpuClusterId)
	}

	if instance.MaintenancePolicy != compute.MaintenancePolicy_MAINTENANCE_POLICY_UNSPECIFIED {
		if err := d.Set("maintenance_policy", strings.ToLower(instance.MaintenancePolicy.String())); err != nil {
			return err
		}
	}

	if err := d.Set("maintenance_grace_period", formatDuration(instance.MaintenanceGracePeriod)); err != nil {
		return err
	}

	if err := d.Set("hardware_generation", hardwareGeneration); err != nil {
		return err
	}

	d.SetId(instance.Id)

	return nil
}
