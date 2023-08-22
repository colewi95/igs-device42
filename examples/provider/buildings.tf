data "device42_buildings" "all" {}

output "all_buildings" {
  value = data.device42_buildings.all
}
