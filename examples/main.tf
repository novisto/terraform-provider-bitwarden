terraform {
  required_providers {
    bitwarden = {
      source = "novisto/bitwarden"
    }
  }
  required_version = ">= 1.0.3"
}

provider "bitwarden" {}

locals {
  notes = {
    "Note 00" = "NOTE 09"
    "Note 01" = "NOTE 1"
    "Note 02" = "NOTE 2"
    "Note 03" = "NOTE 39"
    "Note 04" = "NOTE 4"
    "Note 05" = "NOTE 5"
    "Note 06" = "NOTE 69"
    "Note 07" = "NOTE 7"
    "Note 08" = "NOTE 8"
    "Note 09" = "NOTE 9"
    "Note 10" = "NOTE 99"
    "Note 11" = "NOTE 9"
    "Note 12" = "NOTE 99"
    "Note 13" = "NOTE 9"
    "Note 14" = "NOTE 99"
    "Note 15" = "NOTE 99"
    "Note 16" = "NOTE 9"
    "Note 17" = "NOTE 9"
    "Note 18" = "NOTE 9"
    "Note 19" = "NOTE 99"
    "Note 20" = "NOTE 99"
    "Note 21" = "NOTE 9"
    "Note 22" = "NOTE 9"
    "Note 23" = "NOTE 9"
    "Note 24" = "NOTE 9"
    "Note 25" = "NOTE 99"
  }
}

resource "bitwarden_secure_note" "platform_db_creds_1" {
  for_each = local.notes

  organization_id = "df4736bb-2f70-47ac-98cb-ad7401042241"
  collection_ids  = ["d42f510e-6f45-404a-8a70-ad8d00f6cadf"]
  name            = each.key
  notes           = each.value
}
