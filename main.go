package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"os/user"
)

func execDocker(volume *string) {
	displayEnv := os.Getenv("DISPLAY")
	cngPathEnv := os.Getenv("CNG_PATH")

	var mountedPath *string
	if len(*volume) > 0 {
		mountedPath = volume
	} else {
		mountedPath = &cngPathEnv
	}

	user, err := user.Current()
	var userUid string
	var userGid string
	if err != nil {
		fmt.Printf("[ERROR] Could not get current user, using default user instead. Error: %a", err)
		userUid = "1000"
		userGid = "1000"
	} else {
		userUid = user.Uid
		userGid = user.Gid
	}

	cmd := exec.Command("docker",
		"run", "--name", "contiker", "-it", "--rm",
		"--privileged",
		"--ipc=host",
		"--network", "host",
		"-e", fmt.Sprintf("DISPLAY=%s", displayEnv),
		"-e", "_JAVA_AWT_WM_NONREPARENTING=1",
		"-e", fmt.Sprintf("LOCAL_UID=%s", userUid), "-e", fmt.Sprintf("LOCAL_GID=%s", userGid),
		"-e", "JDK_JAVA_OPTIONS='-Dawt.useSystemAAFontSettings=on -Dswing.aatext=true -Dswing.defaultlaf=com.sun.java.swing.plaf.gtk.GTKLookAndFeel -Dsun.java2d.opengl=true'",
		"-v", "/dev/:/dev/",
		"-v", "/tmp/.X11-unix:/tmp/.X11-unix",
		"--mount", fmt.Sprintf("type=bind,source=%s,destination=/home/user/contiki-ng", *mountedPath),
		"contiker/contiki-ng", "/bin/bash")

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stdout

	errRun := cmd.Run()
	if errRun != nil {
		fmt.Printf("[ERROR] Command did end in error: %a\n", errRun)
	}
}

func execSh() {
	user, err := user.Current()
	var userUid string
	var userGid string
	if err != nil {
		fmt.Printf("[ERROR] Could not get current user, using default user instead. Error: %a", err)
		userUid = "1000"
		userGid = "1000"
	} else {
		userUid = user.Uid
		userGid = user.Gid
	}

	cmd := exec.Command("docker",
		"exec", "--user", fmt.Sprintf("%s:%s", userUid, userGid), "-it", "contiker", "/bin/bash")

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stdout

	errRun := cmd.Run()
	if errRun != nil {
		fmt.Printf("[ERROR] Command did end in error: %a\n", errRun)
	}
}

func execRm() {
	cmd := exec.Command("docker",
		"container", "rm", "-f", "contiker")

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stdout

	errRun := cmd.Run()
	if errRun != nil {
		fmt.Printf("[ERROR] Command did end in error: %a\n", errRun)
	}
}

func main() {
	shPtr := flag.Bool("sh", false, "Activate shell")
	rmPtr := flag.Bool("rm", false, "Remove current docker")
	volumePtr := flag.String("v", "", "Volume to be mounted")
	flag.Parse()

	// flags should be exclusive
	numFlags := 0
	if *shPtr {
		numFlags += 1
	}
	if *rmPtr {
		numFlags += 1
	}
	if len(*volumePtr) > 0 {
		numFlags += 1
	}

	if numFlags > 1 {
		fmt.Printf("[ERROR] Could only handle one flag at a time\n")
		flag.PrintDefaults()
		os.Exit(-1)
	}

	if *shPtr {
		execSh()
	} else if *rmPtr {
		execRm()
	} else {
		execDocker(volumePtr)
	}
}
