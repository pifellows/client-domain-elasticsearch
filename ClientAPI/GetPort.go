package main

import (
	pb "client-domain-elasticsearch/ClientAPI/PortsCommunication"
	"github.com/gin-gonic/gin"
)

type GetPortParameters struct {
	PortId string `uri:"portid" binding:"required"`
}

func GetPort(ctx *gin.Context) {
	ctx.Header("Content-Type", "application/json")
	var params GetPortParameters
	if err := ctx.ShouldBind(&params); err != nil {
		ctx.JSON(400, gin.H{"msg": err.Error()})
		return
	}
	targetPort := pb.PortID{Id: params.PortId}
	port, err := portsClient.GetPort(ctx, &targetPort)
	if err != nil {
		ctx.JSON(400, gin.H{"msg": err.Error()})
	}
	ctx.JSON(200, port)
}
