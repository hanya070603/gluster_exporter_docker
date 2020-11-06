package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"time"

	"github.com/ofesseler/gluster_exporter/structs"
	"github.com/prometheus/common/log"
)

func execGlusterCommand(cmd string) (*bytes.Buffer, error) {
        stdoutBuffer := &bytes.Buffer{}
	glusterExec := exec.Command("/bin/bash","-c", cmd)
	glusterExec.Stdout = stdoutBuffer
	err := glusterExec.Run()
        log.Info("stdoutBuffer:", stdoutBuffer)
        log.Info("glusterExec:", glusterExec)
	if err != nil {
		log.Errorf("tried to execute %v and got error: %v", cmd, err)
		return stdoutBuffer, err
	}
	return stdoutBuffer, nil
}

func execMountCheck() (*bytes.Buffer, error) {
	stdoutBuffer := &bytes.Buffer{}
	mountCmd := exec.Command("mount", "-t", "fuse.glusterfs")

	mountCmd.Stdout = stdoutBuffer

	return stdoutBuffer, mountCmd.Run()
}

func execTouchOnVolumes(mountpoint string) (bool, error) {
	testFileName := fmt.Sprintf("%v/%v_%v", mountpoint, "gluster_mount.test", time.Now())
	_, createErr := os.Create(testFileName)
	if createErr != nil {
		return false, createErr
	}
	removeErr := os.Remove(testFileName)
	if removeErr != nil {
		return false, removeErr
	}
	return true, nil
}

// ExecVolumeInfo executes "gluster volume info" at the local machine and
// returns VolumeInfoXML struct and error
func ExecVolumeInfo() (structs.VolumeInfoXML, error) {
	//args := []string{"docker","exec","-it","glusterfs","/usr/sbin/gluster","volume", "info"}
        args := "docker exec " + containerName + " sh -c \"/usr/sbin/gluster volume info --xml\""
       //args := "docker exec glusterfs sh -c \"/usr/sbin/gluster volume info --xml\""
	bytesBuffer, cmdErr := execGlusterCommand(args)
	if cmdErr != nil {
		return structs.VolumeInfoXML{}, cmdErr
	}
	volumeInfo, err := structs.VolumeInfoXMLUnmarshall(bytesBuffer)
	if err != nil {
		log.Errorf("Something went wrong while unmarshalling xml: %v", err)
		return volumeInfo, err
	}

	return volumeInfo, nil
}

// ExecVolumeList executes "gluster volume info" at the local machine and
// returns VolumeList struct and error
func ExecVolumeList() (structs.VolList, error) {
        //args := []string{"docker","exec","-it","glusterfs","/usr/sbin/gluster","volume", "info"}
        args := "docker exec " + containerName + " sh -c \"/usr/sbin/gluster volume list --xml\""
        bytesBuffer, cmdErr := execGlusterCommand(args)
	if cmdErr != nil {
		return structs.VolList{}, cmdErr
	}
	volumeList, err := structs.VolumeListXMLUnmarshall(bytesBuffer)
	if err != nil {
		log.Errorf("Something went wrong while unmarshalling xml: %v", err)
		return volumeList.VolList, err
	}

	return volumeList.VolList, nil
}

// ExecPeerStatus executes "gluster peer status" at the local machine and
// returns PeerStatus struct and error
func ExecPeerStatus() (structs.PeerStatus, error) {
        //args := []string{"docker","exec","-it","glusterfs","/usr/sbin/gluster","volume", "info"}
        args := "docker exec " + containerName + " sh -c \"/usr/sbin/gluster peer status --xml\""
        bytesBuffer, cmdErr := execGlusterCommand(args)
	if cmdErr != nil {
		return structs.PeerStatus{}, cmdErr
	}
	peerStatus, err := structs.PeerStatusXMLUnmarshall(bytesBuffer)
	if err != nil {
		log.Errorf("Something went wrong while unmarshalling xml: %v", err)
		return peerStatus.PeerStatus, err
	}

	return peerStatus.PeerStatus, nil
}

// ExecVolumeProfileGvInfoCumulative executes "gluster volume {volume] profile info cumulative" at the local machine and
// returns VolumeInfoXML struct and error
func ExecVolumeProfileGvInfoCumulative(volumeName string) (structs.VolProfile, error) {
        //args := []string{"docker","exec","-it","glusterfs","/usr/sbin/gluster","volume", "info"}
        //args := "docker exec -it glusterfs /usr/sbin/gluster volume profile "+volumeName+" info cumulative --xml"
        args := "docker exec " + containerName + " sh -c \"/usr/sbin/gluster volume profile " + volumeName + " info cumulative --xml\""
        bytesBuffer, cmdErr := execGlusterCommand(args)
	if cmdErr != nil {
		return structs.VolProfile{}, cmdErr
	}
	volumeProfile, err := structs.VolumeProfileGvInfoCumulativeXMLUnmarshall(bytesBuffer)
	if err != nil {
		log.Errorf("Something went wrong while unmarshalling xml: %v", err)
		return volumeProfile.VolProfile, err
	}
	return volumeProfile.VolProfile, nil
}

// ExecVolumeStatusAllDetail executes "gluster volume status all detail" at the local machine
// returns VolumeStatusXML struct and error
func ExecVolumeStatusAllDetail() (structs.VolumeStatusXML, error) {
        //args := []string{"docker","exec","-it","glusterfs","/usr/sbin/gluster","volume", "info"}
        args := "docker exec " + containerName + " sh -c \"/usr/sbin/gluster volume status all detail --xml\""
        bytesBuffer, cmdErr := execGlusterCommand(args)
	if cmdErr != nil {
		return structs.VolumeStatusXML{}, cmdErr
	}
	volumeStatus, err := structs.VolumeStatusAllDetailXMLUnmarshall(bytesBuffer)
	if err != nil {
		log.Errorf("Something went wrong while unmarshalling xml: %v", err)
		return volumeStatus, err
	}
	return volumeStatus, nil
}

// ExecVolumeHealInfo executes volume heal info on host system and processes input
// returns (int) number of unsynced files
func ExecVolumeHealInfo(volumeName string) (int, error) {
        //args := []string{"docker","exec","-it","glusterfs","/usr/sbin/gluster","volume", "info"}
        //args := "docker exec -it glusterfs /usr/sbin/gluster volume heal "+volumeName+" info --xml"
        args := "docker exec " + containerName + " sh -c \"/usr/sbin/gluster volume heal " + volumeName + " info --xml\""
        bytesBuffer, cmdErr := execGlusterCommand(args)
	entriesOutOfSync := 0
	if cmdErr != nil {
		return -1, cmdErr
	}
	healInfo, err := structs.VolumeHealInfoXMLUnmarshall(bytesBuffer)
	if err != nil {
		log.Error(err)
		return -1, err
	}

	for _, brick := range healInfo.HealInfo.Bricks.Brick {
		var count int
		var err error
		count, err = strconv.Atoi(brick.NumberOfEntries)
		if err != nil {
			log.Error(err)
			return -1, err
		}
		entriesOutOfSync += count
	}
	return entriesOutOfSync, nil
}

// ExecVolumeQuotaList executes volume quota list on host system and processes input
// returns QuotaList structs and errors
func ExecVolumeQuotaList(volumeName string) (structs.VolumeQuotaXML, error) {
        //args := []string{"docker","exec","-it","glusterfs","/usr/sbin/gluster","volume", "info"}
        //args := "docker exec -it glusterfs /usr/sbin/gluster volume quota "+volumeName+" list --xml"
        args := "docker exec " + containerName + " sh -c \"/usr/sbin/gluster volume quota " + volumeName + " list --xml\""
        bytesBuffer, cmdErr := execGlusterCommand(args)
	if cmdErr != nil {
		return structs.VolumeQuotaXML{}, cmdErr
	}
	volumeQuota, err := structs.VolumeQuotaListXMLUnmarshall(bytesBuffer)
	if err != nil {
		log.Errorf("Something went wrong while unmarshalling xml: %v", err)
		return volumeQuota, err
	}
	return volumeQuota, nil
}
