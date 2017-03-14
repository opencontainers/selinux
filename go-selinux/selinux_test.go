// +build linux,selinux

package selinux

import (
	"os"
	"testing"
)

func TestSetfilecon(t *testing.T) {
	if SelinuxEnabled() {
		tmp := "selinux_test"
		con := "system_u:object_r:bin_t:s0"
		out, _ := os.OpenFile(tmp, os.O_WRONLY|os.O_CREATE, 0)
		out.Close()
		err := Setfilecon(tmp, con)
		if err != nil {
			t.Log("Setfilecon failed")
			t.Fatal(err)
		}
		filecon, err := Getfilecon(tmp)
		if err != nil {
			t.Log("Getfilecon failed")
			t.Fatal(err)
		}
		if con != filecon {
			t.Fatal("Getfilecon failed, returned %s expected %s", filecon, con)
		}

		os.Remove(tmp)
	}
}

func TestSELinux(t *testing.T) {
	var (
		err            error
		plabel, flabel string
	)

	if SelinuxEnabled() {
		t.Log("Enabled")
		plabel, flabel = GetLxcContexts()
		t.Log(plabel)
		t.Log(flabel)
		FreeLxcContexts(plabel)
		plabel, flabel = GetLxcContexts()
		t.Log(plabel)
		t.Log(flabel)
		FreeLxcContexts(plabel)
		t.Log("getenforce ", SelinuxGetEnforce())
		mode := SelinuxGetEnforceMode()
		t.Log("getenforcemode ", mode)

		defer SelinuxSetEnforce(mode)
		if err := SelinuxSetEnforce(Enforcing); err != nil {
			t.Fatalf("enforcing selinux failed: %v", err)
		}
		if err := SelinuxSetEnforce(Permissive); err != nil {
			t.Fatalf("setting selinux mode to permissive failed: %v", err)
		}
		SelinuxSetEnforce(mode)

		pid := os.Getpid()
		t.Logf("PID:%d MCS:%s\n", pid, IntToMcs(pid, 1023))
		err = Setfscreatecon("unconfined_u:unconfined_r:unconfined_t:s0")
		if err == nil {
			t.Log(Getfscreatecon())
		} else {
			t.Log("setfscreatecon failed", err)
			t.Fatal(err)
		}
		err = Setfscreatecon("")
		if err == nil {
			t.Log(Getfscreatecon())
		} else {
			t.Log("setfscreatecon failed", err)
			t.Fatal(err)
		}
		t.Log(Getpidcon(1))
	}
}
