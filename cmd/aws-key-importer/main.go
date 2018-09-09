package main

import (
	"bufio"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"io/ioutil"
	"os"
	"path"
	"strings"
)

const Name = "aws-key-importer"

func main() {
	var dryRun bool

	keyImportCommand := &cobra.Command{
		Use:   "import [name] [public key file] [region]",
		Short: "Imports a public key into selected AWS region",
		Long:  "Imports a public key into selected AWS region",

		Run: func(cmd *cobra.Command, args []string) {
			var keyName string
			var publicKeyName string
			pubKeyDefault := useDefaultIdRsaPub()
			var awsregion string = "us-east-1"

			switch len(args) {
			case 0:
				keyName = prompt("Key Name", "")
				publicKeyName = prompt("Public key", pubKeyDefault)
				awsregion = prompt("AWS Region", "us-east-1")
				fmt.Println("")
			case 1:
				keyName = args[0]
				publicKeyName = prompt("Public key", pubKeyDefault)
				awsregion = prompt("AWS Region", "us-east-1")
				fmt.Println("")
			case 2:
				keyName = args[0]
				publicKeyName = args[1]
				awsregion = prompt("AWS Region", "us-east-1")
				fmt.Println("")
			default:
				keyName = strings.TrimSpace(args[0])
				publicKeyName = strings.TrimSpace(args[1])
				awsregion = strings.TrimSpace(args[2])
			}

			if keyName == "" {
				fmt.Print("Key keyName is required.\n\n")
				cmd.Usage()
				return
			}

			if publicKeyName == "nil" {
				fmt.Print("Public key file is required.\n\n")
				cmd.Usage()
				return
			}

			if awsregion == "" { //todo: Check presence in regions(0
				fmt.Print("Region is required.\n\n")
				cmd.Usage()
				return
			}

			err := importKeyPair(keyName, publicKeyName, awsregion, dryRun)
			if err != nil {
				fmt.Printf("Could not import key pair: %v\n", err)
			}
		},
	}

	rootCmd := &cobra.Command{Use: Name}
	rootCmd.PersistentFlags().BoolVarP(&dryRun, "dry-run", "", false, "Validates basic readiness (aws creds set and so on)")
	rootCmd.AddCommand(keyImportCommand)
	rootCmd.Execute()
}

func importKeyPair(keyName string, pubKey string, region string, dryRun bool) error {
	data, err := ioutil.ReadFile(pubKey)
	if err != nil {
		return err
	}

	sess := session.Must(session.NewSession())

	label := fmt.Sprintf("%s:", region)
	client := ec2.New(sess, &aws.Config{Region: aws.String(region)})

	if keyPairExists(client, keyName, dryRun) {
		fmt.Print("Keypair with specified name already exists.\n\n")
 		return nil
	}

	regions := regions(client)
	if !Contains(regions, region) {
		fmt.Print("Region with specified name already exists.\n\n")
		return nil
	}

	input := &ec2.ImportKeyPairInput{
		KeyName:           aws.String(keyName),
		PublicKeyMaterial: data,
		DryRun:            aws.Bool(dryRun),
	}

	resp, err := client.ImportKeyPair(input)

	if err != nil {
		errMsg := err.Error()
		switch {
		case dryRun && strings.HasPrefix(errMsg, "DryRunOperation"):
			fmt.Printf("[Dry Run] %-16s keypair '%s' imported\n", label, keyName)
		case strings.HasPrefix(errMsg, "InvalidKeyPair.Duplicate"):
			fmt.Printf("%-16s Seems keypair '%s' already exists.\n", label, keyName)
		default:
			fmt.Printf("%-16s Unexpected error during importing '%s' - %v\n", label, keyName, err)
		}
	} else {
		fmt.Printf("%-16s Imported keypair '%s' - %v\n", label, keyName, *resp.KeyFingerprint)
	}

	return nil
}

func keyPairExists(svc *ec2.EC2, keyName string, dryRun bool) bool {
	input := &ec2.DescribeKeyPairsInput{
		DryRun:   aws.Bool(dryRun),
		KeyNames: []*string{aws.String(keyName)},
	}

	resp, err := svc.DescribeKeyPairs(input)
	if err != nil {
		return false
	}

	return len(resp.KeyPairs) > 0
}

func Contains(a []string, x string) bool {
	for _, n := range a {
		if x == n {
			return true
		}
	}
	return false
}

func regions(svc *ec2.EC2) []string {

	resp, err := svc.DescribeRegions(nil)
	if err != nil {
		panic(err)
	}

	regions := make([]string, len(resp.Regions))
	for i, region := range resp.Regions {
		regions[i] = *region.RegionName
	}

	return regions
}


func prompt(name string, defaultVal string) string {
	var p string

	r := bufio.NewReader(os.Stdin)

	if defaultVal != "" {
		p = fmt.Sprintf("%s [%s]: ", name, defaultVal)
	} else {
		p = fmt.Sprintf("%s: ", name)
	}

	fmt.Print(p)

	val, err := r.ReadString('\n')
	if err != nil {
		panic(err)
	}

	val = strings.TrimSpace(val)

	if val == "" && defaultVal != "" {
		val = defaultVal
	}

	return val
}

func useDefaultIdRsaPub() string {
	var pubKey string
	home, err := homedir.Dir()
	if err != nil {
		pubKey = ""
	} else {
		pubKey = path.Join(home, ".ssh/id_rsa.pub")
	}
	return pubKey
}
