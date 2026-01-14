terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.30"
    }
  }
}

# 核心配置：所有请求转发给 LocalStack
provider "aws" {
  access_key                  = "test"     # 随便填
  secret_key                  = "test"     # 随便填
  region                      = "us-east-1"
  
  # 关键：跳过真实的验证
  skip_credentials_validation = true
  skip_metadata_api_check     = true
  skip_requesting_account_id  = true

  # 关键：将 endpoint 指向本地
  endpoints {
    ec2 = "http://localhost:4566"
    iam = "http://localhost:4566"
    sts = "http://localhost:4566"
    s3  = "http://localhost:4566"
  }
# 🚀 必须加上这些，否则会无限卡死
}

# --- 测试 EC2 ---
resource "aws_instance" "test_vm" {
  ami           = "ami-df5de72ade3b4233" # LocalStack 会忽略这个ID，或者你可以配置它映射到具体镜像
  instance_type = "m5.large"

  tags = {
    Name = "Local-VM-01"
  }
}


resource "aws_instance" "jenkins_vm" {
  ami           = "ami-df5de72ade3b4233" # LocalStack 会忽略这个ID，或者你可以配置它映射到具体镜像
  instance_type = "m5.large"

  tags = {
    Name = "jenkins-vm-01"
  }
}
