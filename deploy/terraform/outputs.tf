output "public_ip" {
  description = "Public IP of the portfolio server"
  value       = aws_eip.portfolio.public_ip
}

output "ssh_command" {
  description = "Command to view the portfolio"
  value       = "ssh ${aws_eip.portfolio.public_ip}"
}

output "admin_ssh_command" {
  description = "Command to SSH into the server for admin"
  value       = "ssh -p 2222 ubuntu@${aws_eip.portfolio.public_ip}"
}

output "instance_id" {
  description = "EC2 instance ID"
  value       = aws_instance.portfolio.id
}
