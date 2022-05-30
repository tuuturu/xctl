package argocd

import (
	"bytes"
	"context"
	"crypto/ed25519"
	"crypto/rand"
	_ "embed"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"io"
	"net/http"
	"strings"
	"text/template"

	"github.com/deifyed/xctl/pkg/tools/secrets"
	"golang.org/x/oauth2"

	"github.com/google/go-github/v44/github"

	"github.com/deifyed/xctl/pkg/config"
	"github.com/mikesmitty/edkey"
	"golang.org/x/crypto/ssh"
)

func generateKey() (keyPair, error) {
	pubKey, privKey, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		return keyPair{}, fmt.Errorf("generating key: %w", err)
	}

	publicKey, err := ssh.NewPublicKey(pubKey)
	if err != nil {
		return keyPair{}, fmt.Errorf("serializing public key: %w", err)
	}

	pemKey := &pem.Block{
		Type:  "OPENSSH PRIVATE KEY",
		Bytes: edkey.MarshalED25519PrivateKey(privKey),
	}

	return keyPair{
		PublicKey:  ssh.MarshalAuthorizedKey(publicKey),
		PrivateKey: pem.EncodeToMemory(pemKey),
	}, nil
}

func generateRepositorySecret(repo repository, privateKey []byte) (io.Reader, error) {
	t, err := template.New("secret").Parse(repositorySecretTemplate)
	if err != nil {
		return nil, fmt.Errorf("parsing template: %w", err)
	}

	buf := bytes.Buffer{}

	err = t.Execute(&buf, repositorySecretOpts{
		SecretName:           toRepositorySecretName(repo.Name()),
		OperationsNamespace:  config.DefaultOperationsNamespace,
		RepositoryName:       b64(repo.Name()),
		RepositoryURI:        b64(repo.URL),
		RepositoryPrivateKey: b64(string(privateKey)),
	})
	if err != nil {
		return nil, fmt.Errorf("generating secret: %w", err)
	}

	return &buf, nil
}

func deleteKey(ctx context.Context, secretClient secrets.Client, clusterName string, repo repository) error {
	client, err := authenticatedClient(secretClient)
	if err != nil {
		return fmt.Errorf("preparing client: %w", err)
	}

	keys, _, err := client.Repositories.ListKeys(ctx, repo.Owner(), repo.Name(), &github.ListOptions{})
	if err != nil {
		return fmt.Errorf("listing keys: %w", err)
	}

	expectedName := deployKeyName(clusterName)
	var keyToRemove int64 = -1

	for _, key := range keys {
		if strings.EqualFold(expectedName, *key.Title) {
			keyToRemove = *key.ID

			break
		}
	}

	if keyToRemove == -1 {
		return nil
	}

	_, err = client.Repositories.DeleteKey(ctx, repo.Owner(), repo.Name(), keyToRemove)
	if err != nil {
		return fmt.Errorf("deleting: %w", err)
	}

	return nil
}

type installDeployKeyOpts struct {
	SecretClient secrets.Client
	ClusterName  string
	Repository   repository
	PublicKey    []byte
}

func installDeployKey(ctx context.Context, opts installDeployKeyOpts) error {
	client, err := authenticatedClient(opts.SecretClient)
	if err != nil {
		return fmt.Errorf("preparing client: %w", err)
	}

	_, _, err = client.Repositories.CreateKey(ctx, opts.Repository.Owner(), opts.Repository.Name(), &github.Key{
		Key:      github.String(string(opts.PublicKey)),
		Title:    github.String(deployKeyName(opts.ClusterName)),
		ReadOnly: github.Bool(true),
		//RESEARCH: Verified flag
	})
	if err != nil {
		return fmt.Errorf("creating: %w", err)
	}

	return nil
}

//go:embed templates/ssh-secret.yaml
var repositorySecretTemplate string

func authenticatedClient(secretClient secrets.Client) (*github.Client, error) {
	accessToken, err := secretClient.Get(config.DefaultSecretsGithubNamespace, config.DefaultSecretsGithubAccessTokenKey)
	if err != nil {
		return nil, fmt.Errorf("retrieving access token: %w", err)
	}

	tokenSource := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: accessToken})

	oauth2Client := &http.Client{
		Transport: &oauth2.Transport{Source: tokenSource},
	}

	return github.NewClient(oauth2Client), nil
}

func deployKeyName(clusterName string) string {
	return fmt.Sprintf("%s-%s-%s", config.ApplicationName, clusterName, strings.ToLower(pluginName))
}

func toRepositorySecretName(name string) string {
	return fmt.Sprintf("xctl-argocd-repository-%s", name)
}

func b64(raw string) string {
	rawAsBytes := []byte(raw)

	result := make([]byte, base64.StdEncoding.EncodedLen(len(rawAsBytes)))

	base64.StdEncoding.Encode(result, rawAsBytes)

	return string(result)
}
