package exec

import (
	"os/exec"
)

// Wrapper of the default `exec.Command` in order to support pre and post hooks.
//
// Example:
// cmd := Command("ls", "-l")
// cmd.Run() # It will execute DEFAULT_PRE_HOOK and DEFAULT_POST_HOOK funcs before and after the command.
//

var DEFAULT_PRE_HOOK, DEFAULT_POST_HOOK func()

type CmdWrapper struct {
	cmd *exec.Cmd
}

func InitDefaultHooks(preHook, postHook func()){
	DEFAULT_PRE_HOOK = preHook
	DEFAULT_POST_HOOK = postHook
}

func Command(name string, arg ...string) *CmdWrapper {
	return &CmdWrapper{
		cmd: exec.Command(name, arg...)}
}

func (cw CmdWrapper) Run() error {
	DEFAULT_PRE_HOOK()
	err:= cw.cmd.Run()
	DEFAULT_POST_HOOK()

	return err
}

func init(){
	DEFAULT_PRE_HOOK = func(){}
	DEFAULT_POST_HOOK = func(){}
}
