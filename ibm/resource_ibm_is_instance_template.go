package ibm

import (
	"fmt"

	"github.com/IBM/vpc-go-sdk/vpcv1"
	"github.com/hashicorp/terraform-plugin-sdk/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

const (
	isInstanceTemplateBootVolume                   = "boot_volume"
	isInstanceTemplateVolAttVolAutoDelete          = "auto_delete"
	isInstanceTemplateVolAttVolCapacity            = "capacity"
	isInstanceTemplateVolAttVolIops                = "iops"
	isInstanceTemplateVolAttVolName                = "name"
	isInstanceTemplateVolAttVolBillingTerm         = "billing_term"
	isInstanceTemplateVolAttVolEncryptionKey       = "encryption_key"
	isInstanceTemplateVolAttVolType                = "type"
	isInstanceTemplateVolAttVolProfile             = "profile"
	isInstanceTemplateProvisioning                 = "provisioning"
	isInstanceTemplateProvisioningDone             = "done"
	isInstanceTemplateAvailable                    = "available"
	isInstanceTemplateDeleting                     = "deleting"
	isInstanceTemplateDeleteDone                   = "done"
	isInstanceTemplateFailed                       = "failed"
	isInstanceTemplateBootName                     = "name"
	isInstanceTemplateBootSize                     = "size"
	isInstanceTemplateBootIOPS                     = "iops"
	isInstanceTemplateBootEncryption               = "encryption"
	isInstanceTemplateBootProfile                  = "profile"
	isInstanceTemplateVolumeAttaching              = "attaching"
	isInstanceTemplateVolumeAttached               = "attached"
	isInstanceTemplateVolumeDetaching              = "detaching"
	isInstanceTemplatePlacementTarget              = "placement_target"
	isInstanceTemplateDedicatedHost                = "dedicated_host"
	isInstanceTemplateDedicatedHostGroup           = "dedicated_host_group"
	isInstanceTemplateResourceType                 = "resource_type"
	isInstanceTemplateVolumeDeleteOnInstanceDelete = "delete_volume_on_instance_delete"
)

func resourceIBMISInstanceTemplate() *schema.Resource {
	return &schema.Resource{
		Create:   resourceIBMisInstanceTemplateCreate,
		Read:     resourceIBMisInstanceTemplateRead,
		Update:   resourceIBMisInstanceTemplateUpdate,
		Delete:   resourceIBMisInstanceTemplateDelete,
		Exists:   resourceIBMisInstanceTemplateExists,
		Importer: &schema.ResourceImporter{},

		CustomizeDiff: customdiff.Sequence(
			func(diff *schema.ResourceDiff, v interface{}) error {
				return resourceTagsCustomizeDiff(diff)
			},
		),

		Schema: map[string]*schema.Schema{
			isInstanceTemplateName: {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     false,
				ValidateFunc: validateISName,
				Description:  "Instance Template name",
			},

			isInstanceTemplateVPC: {
				Type:        schema.TypeString,
				ForceNew:    true,
				Required:    true,
				Description: "VPC id",
			},

			isInstanceTemplateZone: {
				Type:        schema.TypeString,
				ForceNew:    true,
				Required:    true,
				Description: "Zone name",
			},

			isInstanceTemplateProfile: {
				Type:        schema.TypeString,
				ForceNew:    true,
				Required:    true,
				Description: "Profile info",
			},

			isInstanceTemplateKeys: {
				Type:             schema.TypeSet,
				Required:         true,
				Elem:             &schema.Schema{Type: schema.TypeString},
				Set:              schema.HashString,
				DiffSuppressFunc: applyOnce,
				Description:      "SSH key Ids for the instance template",
			},

			isInstanceTemplateVolumeAttachments: {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						isInstanceTemplateVolumeDeleteOnInstanceDelete: {
							Type:     schema.TypeBool,
							Required: true,
						},
						"name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"volume": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},

			isInstanceTemplatePrimaryNetworkInterface: {
				Type:        schema.TypeList,
				MinItems:    1,
				MaxItems:    1,
				Required:    true,
				Description: "Primary Network interface info",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						isInstanceTemplateNicName: {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						isInstanceTemplateNicPrimaryIpv4Address: {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						isInstanceTemplateNicSecurityGroups: {
							Type:     schema.TypeSet,
							Optional: true,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
							Set:      schema.HashString,
						},
						isInstanceTemplateNicSubnet: {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
					},
				},
			},

			isInstanceTemplateNetworkInterfaces: {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						isInstanceTemplateNicName: {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						isInstanceTemplateNicPrimaryIpv4Address: {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						isInstanceTemplateNicSecurityGroups: {
							Type:     schema.TypeSet,
							Optional: true,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
							Set:      schema.HashString,
						},
						isInstanceTemplateNicSubnet: {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
					},
				},
			},

			isInstanceTemplateUserData: {
				Type:        schema.TypeString,
				ForceNew:    true,
				Optional:    true,
				Description: "User data given for the instance",
			},

			isInstanceTemplateImage: {
				Type:        schema.TypeString,
				ForceNew:    true,
				Required:    true,
				Description: "image name",
			},

			isInstanceTemplateBootVolume: {
				Type:             schema.TypeList,
				DiffSuppressFunc: applyOnce,
				Optional:         true,
				Computed:         true,
				MaxItems:         1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						isInstanceTemplateBootName: {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						isInstanceTemplateBootEncryption: {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						isInstanceTemplateBootSize: {
							Type:     schema.TypeInt,
							Optional: true,
							Computed: true,
						},
						isInstanceTemplateBootIOPS: {
							Type:     schema.TypeInt,
							Optional: true,
							Computed: true,
						},
						isInstanceTemplateBootProfile: {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						isInstanceTemplateVolumeDeleteOnInstanceDelete: {
							Type:     schema.TypeBool,
							Optional: true,
							Computed: true,
						},
					},
				},
			},

			isInstanceTemplateResourceGroup: {
				Type:        schema.TypeString,
				ForceNew:    true,
				Optional:    true,
				Computed:    true,
				Description: "Instance template resource group",
			},
		},
	}
}

func resourceIBMisInstanceTemplateCreate(d *schema.ResourceData, meta interface{}) error {
	profile := d.Get(isInstanceTemplateProfile).(string)
	name := d.Get(isInstanceTemplateName).(string)
	vpcID := d.Get(isInstanceTemplateVPC).(string)
	zone := d.Get(isInstanceTemplateZone).(string)
	image := d.Get(isInstanceTemplateImage).(string)

	err := instanceTemplateCreate(d, meta, profile, name, vpcID, zone, image)
	if err != nil {
		return err
	}

	return resourceIBMisInstanceTemplateRead(d, meta)
}

func resourceIBMisInstanceTemplateRead(d *schema.ResourceData, meta interface{}) error {
	ID := d.Id()
	err := instanceTemplateGet(d, meta, ID)
	if err != nil {
		return err
	}
	return nil
}

func resourceIBMisInstanceTemplateDelete(d *schema.ResourceData, meta interface{}) error {

	ID := d.Id()

	err := instanceTemplateDelete(d, meta, ID)
	if err != nil {
		return err
	}
	return nil
}

func resourceIBMisInstanceTemplateUpdate(d *schema.ResourceData, meta interface{}) error {

	err := instanceTemplateUpdate(d, meta)
	if err != nil {
		return err
	}
	return resourceIBMisInstanceTemplateRead(d, meta)
}

func resourceIBMisInstanceTemplateExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	ID := d.Id()
	ok, err := instanceTemplateExists(d, meta, ID)
	if err != nil {
		return false, err
	}
	return ok, err
}

func instanceTemplateCreate(d *schema.ResourceData, meta interface{}, profile, name, vpcID, zone, image string) error {
	sess, err := vpcClient(meta)
	if err != nil {
		return err
	}
	instanceproto := &vpcv1.InstanceTemplatePrototype{
		Image: &vpcv1.ImageIdentity{
			ID: &image,
		},
		Zone: &vpcv1.ZoneIdentity{
			Name: &zone,
		},
		Profile: &vpcv1.InstanceProfileIdentity{
			Name: &profile,
		},
		Name: &name,
		VPC: &vpcv1.VPCIdentity{
			ID: &vpcID,
		},
	}

	// BOOT VOLUME ATTACHMENT for instance template
	if boot, ok := d.GetOk(isInstanceTemplateBootVolume); ok {
		bootvol := boot.([]interface{})[0].(map[string]interface{})
		var volTemplate = &vpcv1.VolumePrototypeInstanceByImageContext{}
		name, ok := bootvol[isInstanceTemplateBootName]
		namestr := name.(string)
		if ok {
			volTemplate.Name = &namestr
		}

		if volcap, ok := bootvol[isInstanceTemplateBootSize]; ok {
			if volcapint64 := int64(volcap.(int)); volcapint64 != 0 {
				volTemplate.Capacity = &volcapint64
			}
		}

		if volprof, ok := bootvol[isInstanceTemplateBootProfile]; ok {
			if volumeProfile := volprof.(string); volumeProfile != "" {
				volTemplate.Profile = &vpcv1.VolumeProfileIdentity{
					Name: &volumeProfile,
				}
			}
		}

		if iops, ok := bootvol[isInstanceTemplateBootIOPS]; ok {
			if bootVolIOPS := int64(iops.(int)); bootVolIOPS != 0 {
				volTemplate.Iops = &bootVolIOPS
			}
		}

		if encryption, ok := bootvol[isInstanceTemplateBootEncryption]; ok {
			bootEncryption := encryption.(string)
			if bootEncryption != "" {
				volTemplate.EncryptionKey = &vpcv1.EncryptionKeyIdentity{
					CRN: &bootEncryption,
				}
			}
		}

		var deleteVolumeOption bool
		if deleteVolume, ok := bootvol[isInstanceTemplateVolumeDeleteOnInstanceDelete]; ok {
			deleteVolumeOption = deleteVolume.(bool)
		}

		instanceproto.BootVolumeAttachment = &vpcv1.VolumeAttachmentPrototypeInstanceByImageContext{
			DeleteVolumeOnInstanceDelete: &deleteVolumeOption,
			Volume:                       volTemplate,
		}
	}

	// Handle volume attachments
	if volsintf, ok := d.GetOk(isInstanceTemplateVolumeAttachments); ok {
		vols := volsintf.([]interface{})
		var intfs []vpcv1.VolumeAttachmentPrototypeInstanceContext
		for _, resource := range vols {
			vol := resource.(map[string]interface{})
			volInterface := &vpcv1.VolumeAttachmentPrototypeInstanceContext{}
			deleteVol, _ := vol[isInstanceTemplateVolumeDeleteOnInstanceDelete]
			deleteVolBool := deleteVol.(bool)
			volInterface.DeleteVolumeOnInstanceDelete = &deleteVolBool
			name, _ := vol["name"]
			namestr := name.(string)
			volInterface.Name = &namestr
			volintf, _ := vol["volume"]
			volintfstr := volintf.(string)
			volInterface.Volume = &vpcv1.VolumeAttachmentPrototypeInstanceContextVolume{
				ID: &volintfstr,
			}
			intfs = append(intfs, *volInterface)
		}
		instanceproto.VolumeAttachments = intfs
	}

	// Handle primary network interface
	if primnicintf, ok := d.GetOk(isInstanceTemplatePrimaryNetworkInterface); ok {
		primnic := primnicintf.([]interface{})[0].(map[string]interface{})
		subnetintf, _ := primnic[isInstanceTemplateNicSubnet]
		subnetintfstr := subnetintf.(string)
		var primnicobj = &vpcv1.NetworkInterfacePrototype{}
		primnicobj.Subnet = &vpcv1.SubnetIdentity{
			ID: &subnetintfstr,
		}

		if name, ok := primnic[isInstanceTemplateNicName]; ok {
			namestr := name.(string)
			if namestr != "" {
				primnicobj.Name = &namestr
			}
		}

		secgrpintf, ok := primnic[isInstanceTemplateNicSecurityGroups]
		if ok {
			secgrpSet := secgrpintf.(*schema.Set)
			if secgrpSet.Len() != 0 {
				var secgrpobjs = make([]vpcv1.SecurityGroupIdentityIntf, secgrpSet.Len())
				for i, secgrpIntf := range secgrpSet.List() {
					secgrpIntfstr := secgrpIntf.(string)
					secgrpobjs[i] = &vpcv1.SecurityGroupIdentity{
						ID: &secgrpIntfstr,
					}
				}
				primnicobj.SecurityGroups = secgrpobjs
			}
		}
		instanceproto.PrimaryNetworkInterface = primnicobj

		if IPAddress, ok := primnic[isInstanceTemplateNicPrimaryIpv4Address]; ok {
			if PrimaryIpv4Address := IPAddress.(string); PrimaryIpv4Address != "" {
				primnicobj.PrimaryIpv4Address = &PrimaryIpv4Address
			}
		}
	}

	// Handle  additional network interface
	if nicsintf, ok := d.GetOk(isInstanceTemplateNetworkInterfaces); ok {
		nics := nicsintf.([]interface{})
		var intfs []vpcv1.NetworkInterfacePrototype
		for _, resource := range nics {
			nic := resource.(map[string]interface{})
			nwInterface := &vpcv1.NetworkInterfacePrototype{}
			subnetintf, _ := nic[isInstanceTemplateNicSubnet]
			subnetintfstr := subnetintf.(string)
			nwInterface.Subnet = &vpcv1.SubnetIdentity{
				ID: &subnetintfstr,
			}

			name, ok := nic[isInstanceTemplateNicName]
			namestr := name.(string)
			if ok && namestr != "" {
				nwInterface.Name = &namestr
			}

			secgrpintf, ok := nic[isInstanceTemplateNicSecurityGroups]
			if ok {
				secgrpSet := secgrpintf.(*schema.Set)
				if secgrpSet.Len() != 0 {
					var secgrpobjs = make([]vpcv1.SecurityGroupIdentityIntf, secgrpSet.Len())
					for i, secgrpIntf := range secgrpSet.List() {
						secgrpIntfstr := secgrpIntf.(string)
						secgrpobjs[i] = &vpcv1.SecurityGroupIdentity{
							ID: &secgrpIntfstr,
						}
					}
					nwInterface.SecurityGroups = secgrpobjs
				}
			}
			if IPAddress, ok := nic[isInstanceTemplateNicPrimaryIpv4Address]; ok {
				if PrimaryIpv4Address := IPAddress.(string); PrimaryIpv4Address != "" {
					nwInterface.PrimaryIpv4Address = &PrimaryIpv4Address
				}
			}
			intfs = append(intfs, *nwInterface)
		}
		instanceproto.NetworkInterfaces = intfs
	}

	// Handle SSH Keys
	keySet := d.Get(isInstanceTemplateKeys).(*schema.Set)
	if keySet.Len() != 0 {
		keyobjs := make([]vpcv1.KeyIdentityIntf, keySet.Len())
		for i, key := range keySet.List() {
			keystr := key.(string)
			keyobjs[i] = &vpcv1.KeyIdentity{
				ID: &keystr,
			}
		}
		instanceproto.Keys = keyobjs
	}

	// Handle user data
	if userdata, ok := d.GetOk(isInstanceTemplateUserData); ok {
		userdatastr := userdata.(string)
		instanceproto.UserData = &userdatastr
	}

	// handle resource group
	if grp, ok := d.GetOk(isInstanceTemplateResourceGroup); ok {
		grpstr := grp.(string)
		instanceproto.ResourceGroup = &vpcv1.ResourceGroupIdentity{
			ID: &grpstr,
		}

	}

	options := &vpcv1.CreateInstanceTemplateOptions{
		InstanceTemplatePrototype: instanceproto,
	}

	instanceIntf, response, err := sess.CreateInstanceTemplate(options)
	if err != nil {
		return fmt.Errorf("Error creating InstanceTemplate: %s\n%s", err, response)
	}
	instance := instanceIntf.(*vpcv1.InstanceTemplate)
	d.SetId(*instance.ID)
	return nil
}

func instanceTemplateGet(d *schema.ResourceData, meta interface{}, ID string) error {
	instanceC, err := vpcClient(meta)
	if err != nil {
		return err
	}
	getinsOptions := &vpcv1.GetInstanceTemplateOptions{
		ID: &ID,
	}
	instanceIntf, response, err := instanceC.GetInstanceTemplate(getinsOptions)
	if err != nil {
		return fmt.Errorf("Error Getting Instance template: %s\n%s", err, response)
	}
	instance := instanceIntf.(*vpcv1.InstanceTemplate)
	d.Set(isInstanceTemplateName, *instance.Name)
	if instance.Profile != nil {
		instanceProfileIntf := instance.Profile
		identity := instanceProfileIntf.(*vpcv1.InstanceProfileIdentity)
		d.Set(isInstanceTemplateProfile, *identity.Name)
	}

	if instance.PrimaryNetworkInterface != nil {
		primaryNicList := make([]map[string]interface{}, 0)
		currentPrimNic := map[string]interface{}{}
		currentPrimNic[isInstanceTemplateNicName] = *instance.PrimaryNetworkInterface.Name
		if instance.PrimaryNetworkInterface.PrimaryIpv4Address != nil {
			currentPrimNic[isInstanceTemplateNicPrimaryIpv4Address] = *instance.PrimaryNetworkInterface.PrimaryIpv4Address
		}
		subInf := instance.PrimaryNetworkInterface.Subnet
		subnetIdentity := subInf.(*vpcv1.SubnetIdentity)
		currentPrimNic[isInstanceTemplateNicSubnet] = *subnetIdentity.ID

		if len(instance.PrimaryNetworkInterface.SecurityGroups) != 0 {
			secgrpList := []string{}
			for i := 0; i < len(instance.PrimaryNetworkInterface.SecurityGroups); i++ {
				secGrpInf := instance.PrimaryNetworkInterface.SecurityGroups[i]
				subnetIdentity := secGrpInf.(*vpcv1.SecurityGroupIdentity)
				secgrpList = append(secgrpList, string(*subnetIdentity.ID))
			}
			currentPrimNic[isInstanceTemplateNicSecurityGroups] = newStringSet(schema.HashString, secgrpList)
		}
		primaryNicList = append(primaryNicList, currentPrimNic)
		d.Set(isInstanceTemplatePrimaryNetworkInterface, primaryNicList)
	}

	if instance.NetworkInterfaces != nil {
		interfacesList := make([]map[string]interface{}, 0)
		for _, intfc := range instance.NetworkInterfaces {
			currentNic := map[string]interface{}{}
			currentNic[isInstanceTemplateNicName] = *intfc.Name
			if intfc.PrimaryIpv4Address != nil {
				currentNic[isInstanceTemplateNicPrimaryIpv4Address] = *intfc.PrimaryIpv4Address
			}

			subInf := intfc.Subnet
			subnetIdentity := subInf.(*vpcv1.SubnetIdentity)
			currentNic[isInstanceTemplateNicSubnet] = *subnetIdentity.ID
			if len(intfc.SecurityGroups) != 0 {
				secgrpList := []string{}
				for i := 0; i < len(intfc.SecurityGroups); i++ {
					secGrpInf := intfc.SecurityGroups[i]
					subnetIdentity := secGrpInf.(*vpcv1.SecurityGroupIdentity)
					secgrpList = append(secgrpList, string(*subnetIdentity.ID))
				}
				currentNic[isInstanceTemplateNicSecurityGroups] = newStringSet(schema.HashString, secgrpList)
			}
			interfacesList = append(interfacesList, currentNic)
		}
		d.Set(isInstanceTemplateNetworkInterfaces, interfacesList)
	}

	if instance.Image != nil {
		imageInf := instance.Image
		imageIdentity := imageInf.(*vpcv1.ImageIdentity)
		d.Set(isInstanceTemplateImage, *imageIdentity.ID)
	}
	vpcInf := instance.VPC
	vpcRef := vpcInf.(*vpcv1.VPCIdentity)
	d.Set(isInstanceTemplateVPC, vpcRef.ID)
	zoneInf := instance.Zone
	zone := zoneInf.(*vpcv1.ZoneIdentity)
	d.Set(isInstanceTemplateZone, *zone.Name)

	interfacesList := make([]map[string]interface{}, 0)
	if instance.VolumeAttachments != nil {
		for _, volume := range instance.VolumeAttachments {
			volumeAttach := map[string]interface{}{}
			volumeAttach[isInstanceTemplateVolAttName] = *volume.Name
			volumeAttach[isInstanceTemplateDeleteVolume] = *volume.DeleteVolumeOnInstanceDelete
			volumeID := map[string]interface{}{}
			volumeIntf := volume.Volume
			volumeInst := volumeIntf.(*vpcv1.VolumeAttachmentPrototypeInstanceContextVolume)
			if volumeInst.Name != nil {
				volumeID["name"] = *volumeInst.Name
			}
			if volumeInst.Iops != nil {
				volumeID["iops"] = *volumeInst.Iops
			}
			volumeAttach[isInstanceTemplateVolAttVolume] = volumeID
			interfacesList = append(interfacesList, volumeAttach)
		}
		d.Set(isInstanceTemplateVolumeAttachments, interfacesList)
	}
	if instance.BootVolumeAttachment != nil {
		bootVolList := make([]map[string]interface{}, 0)
		bootVol := map[string]interface{}{}
		bootVol[isInstanceTemplateBootName] = *instance.BootVolumeAttachment.Name
		bootVol[isInstanceTemplateVolAttVolume] = *instance.BootVolumeAttachment.Volume
		bootVol[isInstanceTemplateDeleteVolume] = *instance.BootVolumeAttachment.DeleteVolumeOnInstanceDelete
		bootVolList = append(bootVolList, bootVol)
		d.Set(isInstanceTemplateBootVolume, bootVolList)
	}

	if instance.ResourceGroup != nil {
		d.Set(isInstanceTemplateResourceGroup, instance.ResourceGroup.ID)
	}
	return nil
}

func instanceTemplateUpdate(d *schema.ResourceData, meta interface{}) error {
	instanceC, err := vpcClient(meta)
	if err != nil {
		return err
	}
	ID := d.Id()

	if d.HasChange(isInstanceName) {
		name := d.Get(isInstanceTemplateName).(string)
		updnetoptions := &vpcv1.UpdateInstanceTemplateOptions{
			ID: &ID,
		}

		instanceTemplatePatchModel := &vpcv1.InstanceTemplatePatch{
			Name: &name,
		}
		instanceTemplatePatch, err := instanceTemplatePatchModel.AsPatch()
		if err != nil {
			return fmt.Errorf("Error calling asPatch for InstanceTemplatePatch: %s", err)
		}
		updnetoptions.InstanceTemplatePatch = instanceTemplatePatch

		_, _, err = instanceC.UpdateInstanceTemplate(updnetoptions)
		if err != nil {
			return err
		}
	}
	return nil
}

func instanceTemplateDelete(d *schema.ResourceData, meta interface{}, ID string) error {
	instanceC, err := vpcClient(meta)
	if err != nil {
		return err
	}

	deleteinstanceTemplateOptions := &vpcv1.DeleteInstanceTemplateOptions{
		ID: &ID,
	}
	_, err = instanceC.DeleteInstanceTemplate(deleteinstanceTemplateOptions)
	if err != nil {
		return err
	}
	return nil
}

func instanceTemplateExists(d *schema.ResourceData, meta interface{}, ID string) (bool, error) {
	instanceC, err := vpcClient(meta)
	if err != nil {
		return false, err
	}
	getinsOptions := &vpcv1.GetInstanceTemplateOptions{
		ID: &ID,
	}
	_, response, err := instanceC.GetInstanceTemplate(getinsOptions)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			return false, nil
		}
		return false, fmt.Errorf("Error Getting InstanceTemplate: %s\n%s", err, response)
	}
	return true, nil
}
