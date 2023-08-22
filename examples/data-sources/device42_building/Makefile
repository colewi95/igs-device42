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
