terraform {
  required_providers {
    bitwarden = {
      source = "registry.novisto.net/novisto/bitwarden"

    }
  }
  required_version = ">= 1.0.3"
}

provider "bitwarden" {}

locals {
  platform_db_creds = {
    novisto : {
      host : "host.example.com"
      port : 5432
      username : "user"
      password : "df765287b64e51"
    }
  }
}

resource "bitwarden_secure_note" "platform_db_creds" {
  organization_id = "df4736bb-2f70-47ac-98cb-ad7401042241"
  collection_ids  = ["d42f510e-6f45-404a-8a70-ad8d00f6cadf"]
  name            = "TEST DELETE ME - Platform DB Credentials"
  notes           = jsonencode(local.platform_db_creds)
}
