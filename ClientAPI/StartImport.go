package main

import (
	pb "client-domain-elasticsearch/ClientAPI/PortsCommunication"
	"client-domain-elasticsearch/ClientAPI/json-chunk-reader"
	"github.com/gin-gonic/gin"
	"io"
	"os"
)

func StartImport(ctx *gin.Context) {
	stream, err := portsClient.StreamPorts(ctx)
	if err != nil {
		ctx.JSON(400, gin.H{"msg": err.Error()})
		return
	}

	f, err := os.Open(filePath)
	if err != nil {
		ctx.JSON(400, gin.H{"msg": err.Error()})
		return
	}

	jcr, err := json_chunk_reader.NewJsonChunkReader(f)
	if err != nil {
		ctx.JSON(400, gin.H{"msg": err.Error()})
		return
	}

	for {
		var port pb.Port
		// parse json
		key, err := jcr.ReadItem(&port)
		if err == io.EOF {
			break
		}
		if err != nil {
			ctx.JSON(400, gin.H{"msg": err.Error()})
			return
		}

		port.Id = key.(string)
		if err := stream.Send(&port); err != nil {
			ctx.JSON(400, gin.H{"msg": err.Error()})
			return
		}
	}
	reply, err := stream.CloseAndRecv()
	if err != nil {
		ctx.JSON(400, gin.H{"msg": err.Error()})
		return
	}
	ctx.JSON(200, gin.H{"msg": reply})
	return
}
