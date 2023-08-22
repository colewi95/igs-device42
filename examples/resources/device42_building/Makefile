debug_log_path = /tmp/terraform.provider.$(notdir $(CURDIR)).debug.log

default: init
init: clean
	terragrunt init -upgrade
clean:
	rm -rf .terraform .terraform.lock.hcl
plan:
	terragrunt plan
apply:
	terragrunt apply
output:
	terragrunt output
destroy:
	terragrunt destroy
debug:
	export TF_LOG=debug
	export TF_LOG_PATH=$(debug_log_path)
undebug:
	export TF_LOG=
	epxort TF_LOG_PATH=
