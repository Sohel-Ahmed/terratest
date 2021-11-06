module testing

go 1.15

require (
	cloud.google.com/go v0.65.0
	cloud.google.com/go/storage v1.10.0
	github.com/Azure/azure-sdk-for-go v46.0.0+incompatible
	github.com/Azure/go-autorest/autorest v0.11.5
	github.com/Azure/go-autorest/autorest/azure/auth v0.5.1
	github.com/aws/aws-lambda-go v1.13.3
	github.com/aws/aws-sdk-go v1.38.0
	github.com/ghodss/yaml v1.0.0
	github.com/go-sql-driver/mysql v1.5.0
	github.com/google/go-containerregistry v0.0.0-20200110202235-f4fb41bf00a3
	github.com/google/uuid v1.1.1
	github.com/gruntwork-io/go-commons v0.10.0
	github.com/gruntwork-io/terratest v0.32.9
	github.com/hashicorp/go-multierror v1.1.1
	github.com/hashicorp/go-version v1.3.0
	github.com/hashicorp/hcl/v2 v2.10.1
	github.com/hashicorp/terraform-json v0.13.0
	github.com/jinzhu/copier v0.3.2
	github.com/jstemmer/go-junit-report v0.9.1
	github.com/magiconair/properties v1.8.0
	github.com/mattn/go-zglob v0.0.2-0.20190814121620-e3c945676326
	github.com/miekg/dns v1.1.43
	github.com/mitchellh/go-homedir v1.1.0
	github.com/oracle/oci-go-sdk v7.1.0+incompatible
	github.com/pquerna/otp v1.2.1-0.20191009055518-468c2dd2b58d
	github.com/sirupsen/logrus v1.6.0
	github.com/stretchr/testify v1.6.1
	github.com/urfave/cli v1.22.4
	github.com/zclconf/go-cty v1.9.1
	golang.org/x/crypto v0.0.0-20210317152858-513c2a44f670
	golang.org/x/net v0.0.0-20210226172049-e18ecbb05110
	golang.org/x/oauth2 v0.0.0-20200902213428-5d25da1a8d43
	google.golang.org/api v0.30.0
	google.golang.org/genproto v0.0.0-20200825200019-8632dd797987
	k8s.io/api v0.19.3
	k8s.io/apimachinery v0.19.3
	k8s.io/client-go v0.19.3
)

replace github.com/gruntwork-io/terratest => github.com/ffernandezcast/terratest v0.28.6-0.20210507100311-be0175b196d4
