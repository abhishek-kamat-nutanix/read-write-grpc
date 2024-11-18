package main

import (
	"context"
	"fmt"
	"log"

	v2 "github.com/kubernetes-csi/external-snapshotter/client/v8/apis/volumesnapshot/v1"
	"github.com/kubernetes-csi/external-snapshotter/client/v8/clientset/versioned"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	//batchv1 "k8s.io/api/batch/v1"
	//"time"
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
	log.Printf("ss created %s \n",*ss.Spec.Source.PersistentVolumeClaimName)

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
	log.Printf("pvc created %s\n",create_pvc.Name)
	
	
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

	// if volumeName == "yaml-pv-claim" {

	// 	var completions int32 = 1
	// 	var UID int64 = 0
	// 	labels := make(map[string]string)
	// 	labels["app"]="app-yamls"
	// 	command := []string{"/bin/sh","-c"}
	// 	str := "kubectl apply -f /yaml/manifests.yaml"
	// 	args := []string{str}
	// 	// app-config-job job comes here
	// 	readerjob:= batchv1.Job{TypeMeta: metav1.TypeMeta{Kind: "Job",APIVersion: "batch"},
	// 							ObjectMeta: metav1.ObjectMeta{Name: "app-config-job"},
	// 							Spec: batchv1.JobSpec{Completions: &completions,
	// 								Template: v1.PodTemplateSpec{ObjectMeta:  metav1.ObjectMeta{Labels: labels}, 
	// 								Spec: v1.PodSpec{RestartPolicy: "OnFailure",ImagePullSecrets: []v1.LocalObjectReference{{Name: "my-registry-secret"}}, 
	// 								Volumes: []v1.Volume{{Name: "yamls",VolumeSource: v1.VolumeSource{PersistentVolumeClaim: &v1.PersistentVolumeClaimVolumeSource{ClaimName: "yaml-pv-claim"}}}},
	// 								Containers: []v1.Container{{Name: "yaml-applier", Image: "bitnami/kubectl:latest",Command: command,Args: args ,SecurityContext: &v1.SecurityContext{RunAsUser: &UID}, VolumeMounts: []v1.VolumeMount{{Name: "yamls", MountPath: "/yaml"}}}},}}},}

	// 	reader, err := clientset.BatchV1().Jobs(namespace).Create(context.Background(),&readerjob,metav1.CreateOptions{})
	// 	if err != nil{
	// 		fmt.Printf("error while creating app-config-job in %v namespace: %v\n",namespace,err)
	// 	}
	// 	log.Printf("job created %v\n",reader.Name)

	// 	job, err := clientset.BatchV1().Jobs(namespace).Get(context.Background(),"app-config-job",metav1.GetOptions{})
	// if err != nil{
	// 	fmt.Printf("error while getting job in %v namespace: %v\n",namespace,err)
	// }

	// deletePolicy := metav1.DeletePropagationBackground
	// flag := 0
	// 	for {
	// 		time.Sleep(10 * time.Second)
	// 	for _, condition := range job.Status.Conditions {
	// 		if condition.Type == batchv1.JobComplete && condition.Status == v1.ConditionTrue {
	// 			//delete diskreader job so diskreader pvc is not bound and can be deleted successfully
	// 			err = clientset.BatchV1().Jobs(namespace).Delete(context.Background(),"app-config-job",metav1.DeleteOptions{PropagationPolicy: &deletePolicy})	
	// 			if err != nil{
	// 				fmt.Printf("error while deleting job in %v namespace: %v\n",namespace,err)
	// 			}
	// 			log.Print("job completed, deleting Job and Pod \n")

	// 			//delete pvc now
	// 			err := clientset.CoreV1().PersistentVolumeClaims(namespace).Delete(context.Background(),volumeName,metav1.DeleteOptions{})
	// 			if err != nil{
	// 				fmt.Printf("error while deleting pvc from %v namespace: %v\n",namespace ,err)
	// 			}
	// 			log.Printf("pvc deleted %v \n",volumeName)
	// 			flag=1
	// 			break;
	// 		} 
	// 	}
	// 	if flag==1 {
	// 		break
	// 	}
	// 	job, err = clientset.BatchV1().Jobs(namespace).Get(context.Background(),"app-config-job",metav1.GetOptions{})	
	// 			if err != nil{
	// 				fmt.Printf("error while getting pvc in %v namespace: %v\n",namespace,err)
	// 			}
		
	// 	}

	// }
	
}

