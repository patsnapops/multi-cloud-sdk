package io_test

import (
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/alibabacloud-go/tea/tea"
	"github.com/joho/godotenv"

	"github.com/xops-infra/multi-cloud-sdk/pkg/io"
	"github.com/xops-infra/multi-cloud-sdk/pkg/model"
)

var TencentIo model.CloudIO

func init() {
	err := godotenv.Load(".env")
	if err != nil {
		panic(err)
	}
	profiles := []model.ProfileConfig{
		{
			Name:  "tencent",
			Cloud: model.TENCENT,
			AK:    os.Getenv("TENCENT_ACCESS_KEY"),
			SK:    os.Getenv("TENCENT_SECRET_KEY"),
		}, {
			Name:  "tx-dev",
			Cloud: model.TENCENT,
			AK:    os.Getenv("TX_DEV_ID"),
			SK:    os.Getenv("TX_DEV_KEY"),
		},
	}
	clientIo := io.NewCloudClient(profiles)
	TencentIo = io.NewTencentClient(clientIo)
}

func TestQueryTencentEmrCluster(t *testing.T) {
	timeStart := time.Now()
	filter := model.EmrFilter{
		Profile: tea.String("tencent"),
		Region:  tea.String("na-ashburn"),
		ClusterStates: []model.EMRClusterStatus{
			model.EMRClusterRunning,
		},
	}
	instances, err := TencentIo.QueryEmrCluster(filter)
	if err != nil {
		t.Error(err)
		return
	}
	for _, instance := range instances.Clusters {
		fmt.Println(tea.Prettify(instance))
	}
	t.Log("Success.", time.Since(timeStart), len(instances.Clusters))
}

func TestDescribeTencentEmrCluster(t *testing.T) {
	timeStart := time.Now()
	instances, err := TencentIo.DescribeEmrCluster(model.DescribeInput{
		Profile: tea.String("tencent"),
		Region:  tea.String("ap-shanghai"),
		// IDS:     []*string{tea.String("emr-alhn4h4s")},
	})
	if err != nil {
		t.Error(err)
		return
	}
	for _, instance := range instances {
		fmt.Println(tea.Prettify(instance))
	}
	t.Log("Success.", time.Since(timeStart), len(instances), instances)
}

func TestCreateEmrCluster(t *testing.T) {
	timeStart := time.Now()
	input := model.CreateEmrClusterInput{
		Name: tea.String("test"),
		Tags: model.Tags{
			{
				Key:   "Owner",
				Value: "zhoushoujian",
			},
		},
		InstanceChargeType: model.POSTPAID_BY_HOUR.String(),
		APPs:               []*string{tea.String("Spark")},
		ResourceSpec: &model.ResourceSpec{
			HA:     tea.Bool(false),
			VPC:    tea.String("vpc-gjljk6e8"),
			Subnet: tea.String("subnet-j94dsqaj"),
			SgId:   tea.String("sg-2qt3di24"),
			KeyID:  tea.String("skey-gyqojq9d"),
			MasterResourceSpec: &model.EMRInstaceSpec{
				InstanceCount: tea.Int64(1),
				InstanceType:  tea.String("TS5.2XLARGE32"),
				DiskType:      tea.String("CLOUD_SSD"),
				DiskNum:       tea.Int64(0),
				DiskSize:      tea.Int64(40),
				RootSize:      tea.Int64(40),
			},
			CoreResourceSpec: &model.EMRInstaceSpec{
				InstanceCount: tea.Int64(2),
				InstanceType:  tea.String("TS5.2XLARGE32"),
				DiskType:      tea.String("CLOUD_SSD"),
				DiskNum:       tea.Int64(0),
				DiskSize:      tea.Int64(40),
				RootSize:      tea.Int64(40),
			},
		},
	}
	instances, err := TencentIo.CreateEmrCluster("tencent", "ap-shanghai", input)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log("Success.", time.Since(timeStart), instances)
}

func TestListInstance(t *testing.T) {
	timeStart := time.Now()
	filter := model.InstanceFilter{}
	instances, err := TencentIo.DescribeInstances("tencent", "ap-beijing", filter.ToTxDescribeInstancesInput())
	if err != nil {
		t.Error(err)
		return
	}
	t.Log("Success.", time.Since(timeStart), len(instances.Instances))
}

func TestDescribeInstances(t *testing.T) {
	{
		timeStart := time.Now()
		filter := model.InstanceFilter{}
		filter.PrivateIp = tea.String(os.Getenv("TEST_TENCENT_PRIVATE_IP"))
		instances, err := TencentIo.DescribeInstances("tencent", "ap-shanghai", filter.ToTxDescribeInstancesInput())
		if err != nil {
			t.Error(err)
			return
		}
		t.Log("PrivateIp Success.", time.Since(timeStart), len(instances.Instances))
	}
	{
		timeStart := time.Now()
		filter := model.InstanceFilter{}
		filter.PublicIp = tea.String(os.Getenv("TEST_TENCENT_PUBLIC_IP"))
		instances, err := TencentIo.DescribeInstances("tencent", "ap-shanghai", filter.ToTxDescribeInstancesInput())
		if err != nil {
			t.Error(err)
			return
		}
		t.Log("PublicIp Success.", time.Since(timeStart), len(instances.Instances))
	}
	{
		timeStart := time.Now()
		filter := model.InstanceFilter{}
		filter.Owner = tea.String("zhoushoujian")
		instances, err := TencentIo.DescribeInstances("tencent", "ap-shanghai", filter.ToTxDescribeInstancesInput())
		if err != nil {
			t.Error(err)
			return
		}
		t.Log("Owner Success.", time.Since(timeStart), len(instances.Instances))
	}
	{
		timeStart := time.Now()
		filter := model.InstanceFilter{}
		filter.IDs = []*string{tea.String(os.Getenv("TEST_TENCENT_ID"))}
		instances, err := TencentIo.DescribeInstances("tencent", "ap-shanghai", filter.ToTxDescribeInstancesInput())
		if err != nil {
			t.Error(err)
			return
		}
		t.Log("ID Success.", time.Since(timeStart), len(instances.Instances))
	}
	{
		timeStart := time.Now()
		filter := model.InstanceFilter{}
		filter.Name = tea.String(os.Getenv("TEST_TENCENT_NAME"))
		instances, err := TencentIo.DescribeInstances("tencent", "ap-shanghai", filter.ToTxDescribeInstancesInput())
		if err != nil {
			t.Error(err)
			return
		}
		t.Log("Name Success.", time.Since(timeStart), len(instances.Instances))
	}
	{
		timeStart := time.Now()
		filter := model.InstanceFilter{}
		filter.Status = model.InstanceStatusStopped.TString()
		instances, err := TencentIo.DescribeInstances("tencent", "ap-shanghai", filter.ToTxDescribeInstancesInput())
		if err != nil {
			t.Error(err)
			return
		}
		t.Log("Status Success.", time.Since(timeStart), len(instances.Instances))
	}
}

func TestCreateTags(t *testing.T) {
	input := model.CreateTagsInput{
		Tags: model.Tags{
			{
				Key:   "CreateTime",
				Value: time.Now().Format("2006010215"),
			},
		},
	}
	err := TencentIo.CreateTags("tencent", "ap-shanghai", input)
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("Success. %s", tea.Prettify(input))
}

func TestCreateInstance(t *testing.T) {
	resp, err := TencentIo.CreateInstance("tencent", "ap-shanghai", model.CreateInstanceInput{
		Name:             tea.String("multi-cloud-sdk-test"),
		ImageID:          tea.String("img-hdt9xxkt"),
		InstanceType:     tea.String("SA5.MEDIUM2"),
		KeyIds:           []*string{tea.String(os.Getenv("TEST_TENCENT_KEY_ID"))},
		Zone:             tea.String("ap-shanghai-5"),
		VpcID:            tea.String(os.Getenv("TEST_TENCENT_VPC_ID")),
		SubnetID:         tea.String(os.Getenv("TEST_TENCENT_SUBNET_ID")),
		SecurityGroupIDs: []*string{tea.String(os.Getenv("TEST_TENCENT_SECURITY_GROUP_ID"))},
	})
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("Success. %s", tea.Prettify(resp))
}

func TestModifyInstance(t *testing.T) {
	instancesIds := []*string{tea.String("ins-iwh5ysbx")}

	// // StartInstance
	// {
	// 	resp, err := TencentIo.ModifyInstance("tencent", "ap-shanghai", model.ModifyInstanceInput{
	// 		Action:      model.StartInstance,
	// 		InstanceIDs: instancesIds,
	// 	})
	// 	if err != nil {
	// 		t.Error(err)
	// 		return
	// 	}
	// 	t.Logf("StartInstance Success. %s", tea.Prettify(resp))
	// 	time.Sleep(30 * time.Second)
	// }
	// // RebootInstance
	// {
	// 	resp, err := TencentIo.ModifyInstance("tencent", "ap-shanghai", model.ModifyInstanceInput{
	// 		Action:      model.RebootInstance,
	// 		InstanceIDs: instancesIds,
	// 	})
	// 	if err != nil {
	// 		t.Error(err)
	// 		return
	// 	}
	// 	t.Logf("RebootInstance Success. %s", tea.Prettify(resp))
	// 	time.Sleep(30 * time.Second)
	// }

	// ResetInstance
	{
		resp, err := TencentIo.ModifyInstance("tencent", "ap-shanghai", model.ModifyInstanceInput{
			Action:      model.ResetInstance,
			InstanceIDs: instancesIds,
		})
		if err != nil {
			t.Error(err)
			return
		}
		t.Logf("ResetInstance Success. %s", tea.Prettify(resp))
		time.Sleep(30 * time.Second)
	}

	// StopInstance
	// {
	// 	resp, err := TencentIo.ModifyInstance("tencent", "ap-shanghai", model.ModifyInstanceInput{
	// 		Action:      model.StopInstance,
	// 		InstanceIDs: instancesIds,
	// 	})
	// 	if err != nil {
	// 		t.Error(err)
	// 		return
	// 	}
	// 	t.Logf("StopInstance Success. %s", tea.Prettify(resp))
	// 	time.Sleep(30 * time.Second)
	// }

	// ChangeInstanceType
	// {
	// 	resp, err := TencentIo.ModifyInstance("tencent", "ap-shanghai", model.ModifyInstanceInput{
	// 		Action:       model.ChangeInstanceType,
	// 		InstanceIDs:  instancesIds,
	// 		InstanceType: tea.String("SA5.MEDIUM2"),
	// 	})
	// 	if err != nil {
	// 		t.Error(err)
	// 		return
	// 	}
	// 	t.Logf("ChangeInstanceType Success. %s", tea.Prettify(resp))
	// }
}

func TestChangeInstanceType(t *testing.T) {
	resp, err := TencentIo.ModifyInstance("tencent", "ap-shanghai", model.ModifyInstanceInput{
		Action:       model.ChangeInstanceType,
		InstanceIDs:  []*string{tea.String("ins-k7fdkyi1")},
		InstanceType: tea.String("SA5.2XLARGE32"),
	})
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("Success. %s", tea.Prettify(resp))
}

// TestResetInstance
func TestResetInstance(t *testing.T) {
	resp, err := TencentIo.ModifyInstance("tencent", "ap-shanghai", model.ModifyInstanceInput{
		Action:      model.ResetInstance,
		InstanceIDs: []*string{tea.String("ins-xx")},
	})
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("Success. %s", tea.Prettify(resp))
}

func TestDeleteInstance(t *testing.T) {
	resp, err := TencentIo.DeleteInstance("tencent", "ap-shanghai", model.DeleteInstanceInput{
		InstanceIds: []*string{tea.String("ins-xx")},
	})
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("Success. %s", tea.Prettify(resp))
}

// TEST CreateSecurityGroupWithPolicies
func TestCreateSecurityGroupWithPolicies(t *testing.T) {
	resp, err := TencentIo.CreateSecurityGroupWithPolicies("tencent", "ap-beijing", model.CreateSecurityGroupWithPoliciesInput{
		GroupName:        tea.String("-test"),
		GroupDescription: tea.String("multi-cloud-sdk-test"),
		PolicySet: model.PolicySet{
			Egress: []model.SecurityGroupPolicy{
				{
					Protocol:          tea.String("ALL"),
					Port:              tea.String("ALL"),
					CidrBlock:         tea.String("0.0.0.0/0"),
					PolicyDescription: tea.String("allow all"),
					Action:            tea.String("ACCEPT"),
				},
			},
			Ingress: []model.SecurityGroupPolicy{
				{
					Protocol:          tea.String("ALL"),
					Port:              tea.String("ALL"),
					CidrBlock:         tea.String("0.0.0.0/0"),
					PolicyDescription: tea.String("allow all"),
					Action:            tea.String("ACCEPT"),
				},
			},
		},
	})
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("Success. %s", tea.Prettify(resp))
}

// TEST CreateSecurityGroupPolicies
func TestCreateSecurityGroupPolicies(t *testing.T) {
	allowAll := strings.Split(os.Getenv("TestCreateSecurityGroupPoliciesCidr"), ",")
	ingress := []model.SecurityGroupPolicy{}
	for _, cidr := range allowAll {
		if cidr == "" {
			continue
		}
		ingress = append(ingress, model.SecurityGroupPolicy{
			Protocol:          tea.String("ALL"),
			Port:              tea.String("ALL"),
			CidrBlock:         tea.String(cidr),
			PolicyDescription: tea.String("allow all for office" + "(mcs)"),
			Action:            tea.String("ACCEPT"),
		})
	}
	fmt.Println(tea.Prettify(ingress))
	resp, err := TencentIo.CreateSecurityGroupPolicies("tencent", "ap-beijing", model.CreateSecurityGroupPoliciesInput{
		SecurityGroupId: tea.String("sg-xxx"),
		PolicySet: model.PolicySet{
			Ingress: ingress,
		},
	})
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("Success. %s", tea.Prettify(resp))
}
