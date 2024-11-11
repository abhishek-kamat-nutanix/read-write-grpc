package main

import (
	"context"
	"fmt"

	v2 "github.com/kubernetes-csi/external-snapshotter/client/v8/apis/volumesnapshot/v1"
	"github.com/kubernetes-csi/external-snapshotter/client/v8/clientset/versioned"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

func writer()  {

	volume := "diskwriter-pvc"

	config, err := rest.InClusterConfig()
	if err!= nil {
		fmt.Printf("error getting in-cluster config: %v\n", err)
	}
	
	
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil{
		fmt.Printf("error creating clientset: %v\n", err)
	}
	// get disk-writer pvc details
	pvc, err := clientset.CoreV1().PersistentVolumeClaims(namespace).Get(context.Background(),volume,metav1.GetOptions{})
	if err != nil{
		fmt.Printf("error while getting pvc %v from %v namespace: %v\n", volume,namespace,err)
	}

	size := pvc.Spec.Resources.Requests.Storage()


	clientset2, err := versioned.NewForConfig(config)
	if err !=nil {
		fmt.Printf("error creating clientset: %v\n", err)
	}

	// create snapshot of disk-writer pvc
	//snapClass := "default-snapshotclass"
	snap:= v2.VolumeSnapshot{
		TypeMeta: metav1.TypeMeta{},ObjectMeta: metav1.ObjectMeta{Name: "source-snap"},Spec: v2.VolumeSnapshotSpec{Source: v2.VolumeSnapshotSource{PersistentVolumeClaimName: &volume}, VolumeSnapshotClassName: &snapClass},Status: &v2.VolumeSnapshotStatus{},
	}

	ss, err := clientset2.SnapshotV1().VolumeSnapshots(namespace).Create(context.Background(),&snap,metav1.CreateOptions{})
	if err != nil{
		fmt.Printf("error while creating snapshot of volume %v: %v\n",volume, err)
	}
	fmt.Printf("ss created %s \n",ss.UID)

	// create new pvc with name
	//storageClassName:=  "default-storageclass"
	volumeMode := v1.PersistentVolumeFilesystem
	persistentVolumeAccessMode := v1.ReadWriteOnce
	resourceName:= v1.ResourceStorage
	m := make(v1.ResourceList)
	m[resourceName] = *size
	apiGroup := "snapshot.storage.k8s.io"
	pvclaim := v1.PersistentVolumeClaim{TypeMeta: metav1.TypeMeta{Kind:"PersistentVolumeClaim",APIVersion:"v1"},
										ObjectMeta: metav1.ObjectMeta{Name: volumeName},
										Spec: v1.PersistentVolumeClaimSpec{StorageClassName: &storageClassName, 
											VolumeMode: &volumeMode, 
											Resources: v1.VolumeResourceRequirements{Limits: v1.ResourceList{},Requests: m}, 
											DataSource: &v1.TypedLocalObjectReference{APIGroup: &apiGroup  ,Kind: "VolumeSnapshot" , Name:"source-snap"},
											AccessModes: []v1.PersistentVolumeAccessMode{persistentVolumeAccessMode}},
										Status: v1.PersistentVolumeClaimStatus{}}
										
	

	
	create_pvc, err := clientset.CoreV1().PersistentVolumeClaims(namespace).Create(context.Background(),&pvclaim,metav1.CreateOptions{})
	if err != nil{
		fmt.Printf("error while creating pvc %v in %v namespace: %v\n", volumeName,namespace,err)
	}
	fmt.Printf("pvc created %s\n",create_pvc.UID)
	
	
	for create_pvc.Status.Phase!= v1.ClaimBound {
		create_pvc, err = clientset.CoreV1().PersistentVolumeClaims(namespace).Get(context.Background(),volumeName,metav1.GetOptions{})
		if err != nil{
			fmt.Printf("error while getting pvc in %v namespace: %v\n",namespace,err)
		}
	}

	err = clientset2.SnapshotV1().VolumeSnapshots(namespace).Delete(context.Background(),"source-snap",metav1.DeleteOptions{})
	 if err != nil{
	 fmt.Printf("error while deleting snapshot from %v namespace: %v\n",namespace ,err)
	 }
	
}

