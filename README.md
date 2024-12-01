# GophKeeper

## Working with HTTPS
#### To run the client in HTTPS mode, you need to provide flags or environment variables:

`client -addr "https://localhost" -crypto-crt "*.crt" -crypto-key "*.key" -crypto-ca "*.key"`

where:  
- `-addr` is the server's address  
- `-crypto-ca` is the key used to sign the certificates  
- `-crypto-crt` is the certificate  
- `-crypto-key` is the client's key  

#### To run the server in HTTPS mode, you need to provide flags or environment variables:

`server -addr ":443" -crypto-crt "*.crt" -crypto-key "*.key" -dns "user=postgres password=12345 host=localhost port=5433 dbname=gophkeep"`

where:  
- `-addr` is the server's address  
- `-crypto-crt` is the certificate  
- `-crypto-key` is the server's key  
- `-dns` is the connection string to the database  

## Working with gRPC
#### To run the client in gRPC mode, you need to provide flags or environment variables:

`client -grpc-addr ":3200"`

where:  
- `-grpc-addr` is the server's address  

#### To run the server in gRPC mode, you need to provide flags or environment variables:

`server -grpc-addr ":3200" -dns "user=postgres password=12345 host=localhost port=5433 dbname=gophkeep"`

where:  
- `-grpc-addr` is the server's address  
- `-dns` is the connection string to the database  
