package minio

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

// NewMinioClient cria e retorna um novo cliente para o Minio.
func NewMinioClient() (*minio.Client, error) {
	ctx := context.Background()

	// Endpoint público (o que o navegador vai usar)
	publicEndpoint := os.Getenv("MINIO_PUBLIC_ENDPOINT")
	if publicEndpoint == "" {
		publicEndpoint = "localhost:9000"
	}

	// Endpoint interno (o nome do serviço dentro da rede do Docker)
	internalHost := os.Getenv("MINIO_HOST") + ":9000"

	accessKeyID := os.Getenv("MINIO_ROOT_USER")
	secretAccessKey := os.Getenv("MINIO_ROOT_PASSWORD")
	useSSL := false

	// O SEGREDO ESTÁ AQUI:
	// Criamos um "Transport" customizado para o cliente HTTP do Go.
	// Quando o minio-go tentar se comunicar com o "publicEndpoint" (localhost:9000),
	// nós interceptamos a conexão e redirecionamos o tráfego TCP para o "internalHost" (minio:9000).
	// Isso permite que a biblioteca assine criptograficamente a URL usando "localhost",
	// resolvendo de forma definitiva e correta o erro SignatureDoesNotMatch.
	customTransport := &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
			if addr == publicEndpoint {
				addr = internalHost
			}
			return (&net.Dialer{
				Timeout:   30 * time.Second,
				KeepAlive: 30 * time.Second,
			}).DialContext(ctx, network, addr)
		},
	}

	// Inicializa o cliente Minio passando o endpoint PÚBLICO e o Transport CUSTOMIZADO
	minioClient, err := minio.New(publicEndpoint, &minio.Options{
		Creds:     credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure:    useSSL,
		Transport: customTransport,
	})
	if err != nil {
		return nil, fmt.Errorf("erro ao inicializar o cliente minio: %w", err)
	}

	log.Println("Conexão com o Minio estabelecida com sucesso usando roteamento de rede customizado!")

	// Garante que o bucket 'videos' exista
	bucketName := "videos"
	err = minioClient.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{})
	if err != nil {
		exists, errBucketExists := minioClient.BucketExists(ctx, bucketName)
		if errBucketExists == nil && exists {
			log.Printf("Bucket '%s' já existe.", bucketName)
		} else {
			return nil, fmt.Errorf("erro ao criar ou verificar o bucket '%s': %w", bucketName, err)
		}
	} else {
		log.Printf("Bucket '%s' criado com sucesso.", bucketName)
	}

	return minioClient, nil
}
