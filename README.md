# Terraform Bitwarden Provider


## Running locally

Local setup for development, you will need Go 1.17 and Terraform 1.0.3+

1. Copy the `.terraformrc.example` file into your HOME and change 
   the name to `.terraformrc` and the path inside to your own username
2. From the root of the repo run `make install`
3. In the `examples` directory run `TF_LOG=TRACE t apply -auto-approve`
