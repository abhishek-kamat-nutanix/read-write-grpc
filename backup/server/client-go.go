package main

import (
	"context"
	"fmt"
	"sync"

	v2 "github.com/kubernetes-csi/external-snapshotter/client/v8/apis/volumesnapshot/v1"
	"github.com/kubernetes-csi/external-snapshotter/client/v8/clientset/versioned"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

func writer(volumeName *string, kubeconfig *string)  {


	volume := "diskwriter-pvc"

	config, err := clientcmd.BuildConfigFromFlags("",*kubeconfig)
	
	if err !=nil {
		fmt.Printf("error building config from flags: %s\n",err.Error())
		config, err = rest.InClusterConfig()

		if err!= nil {
			fmt.Printf("error getting kubeconfig: %v", err)
		}
	} 
	
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil{
		fmt.Printf("error creating clientset: %v\n", err)
	}
	// get disk-writer pvc details
	pvc, err := clientset.CoreV1().PersistentVolumeClaims("default").Get(context.Background(),volume,metav1.GetOptions{})
	if err != nil{
		fmt.Printf("error while getting pvc %v from default namespace: %v\n", volume,err)
	}

	size := pvc.Spec.Resources.Requests.Storage()


	clientset2, err := versioned.NewForConfig(config)
	if err !=nil {
		fmt.Printf("error creating clientset: %v\n", err)
	}

	// create snapshot of disk-writer pvc
	snapClass := "default-snapshotclass"
	snap:= v2.VolumeSnapshot{
		TypeMeta: metav1.TypeMeta{},ObjectMeta: metav1.ObjectMeta{Name: "source-snap"},Spec: v2.VolumeSnapshotSpec{Source: v2.VolumeSnapshotSource{PersistentVolumeClaimName: &volume}, VolumeSnapshotClassName: &snapClass},Status: &v2.VolumeSnapshotStatus{},
	}

	ss, err := clientset2.SnapshotV1().VolumeSnapshots("default").Create(context.Background(),&snap,metav1.CreateOptions{})
	if err != nil{
		fmt.Printf("error while creating snapshot of volume %v: %v\n",volume, err)
	}
	fmt.Printf("ss created %s \n",ss.UID)
	// create new pvc with name
	storageClassName:=  "default-storageclass"
	volumeMode := v1.PersistentVolumeFilesystem
	persistentVolumeAccessMode := v1.ReadWriteOnce
	resourceName:= v1.ResourceStorage
	m := make(v1.ResourceList)
	m[resourceName] = *size
	apiGroup := "snapshot.storage.k8s.io"
	pvclaim := v1.PersistentVolumeClaim{TypeMeta: metav1.TypeMeta{Kind:"PersistentVolumeClaim",APIVersion:"v1"},
										ObjectMeta: metav1.ObjectMeta{Name: *volumeName},
										Spec: v1.PersistentVolumeClaimSpec{StorageClassName: &storageClassName, 
											VolumeMode: &volumeMode, 
											Resources: v1.VolumeResourceRequirements{Limits: v1.ResourceList{},Requests: m}, 
											DataSource: &v1.TypedLocalObjectReference{APIGroup: &apiGroup  ,Kind: "VolumeSnapshot" , Name:"source-snap"},
											AccessModes: []v1.PersistentVolumeAccessMode{persistentVolumeAccessMode}},
										Status: v1.PersistentVolumeClaimStatus{}}
										
	

	
	create_pvc, err := clientset.CoreV1().PersistentVolumeClaims("default").Create(context.Background(),&pvclaim,metav1.CreateOptions{})
	if err != nil{
		fmt.Printf("error while creating pvc %v in default namespace: %v\n", *volumeName,err)
	}
	fmt.Printf("pvc created %s\n",create_pvc.UID)
	
	

	for create_pvc.Status.Phase!= v1.ClaimBound {
		create_pvc, err = clientset.CoreV1().PersistentVolumeClaims("default").Get(context.Background(),*volumeName,metav1.GetOptions{})
		if err != nil{
			fmt.Printf("error while getting pvc in default namespace: %v\n",err)
		}
	}

	err = clientset2.SnapshotV1().VolumeSnapshots("default").Delete(context.Background(),"source-snap",metav1.DeleteOptions{})
	 if err != nil{
	 fmt.Printf("error while deleting snapshot from default namespace: %v\n", err)
	 }
	
}

