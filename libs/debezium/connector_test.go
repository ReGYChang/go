package debezium

import (
	"nexdata/pkg/config"
	"testing"
)

func TestSetupOracleConnector(t *testing.T) {
	config.Debezium.Connector.ConnectAddr = "10.13.1.156:8083"
	config.Debezium.Connector.ConnectName = "oracle_conn"

	config.Debezium.Connector.TasksMax = "1"
	config.Debezium.Connector.DBHostname = "10.90.1.207"
	config.Debezium.Connector.DBPort = "1521"
	config.Debezium.Connector.DBServerName = "oracleserver"
	config.Debezium.Connector.DBUser = "logminer"
	config.Debezium.Connector.DBPassword = "logminer"
	config.Debezium.Connector.DBName = "EMESHY"
	config.Debezium.Connector.TableName = "EMESP.TP_SN_LOG"
	config.Debezium.Connector.DBBootstrapServer = "kafka:9092"
	config.Debezium.Connector.ConnectorDBType = "io.debezium.connector.oracle.OracleConnector"
	config.Debezium.Connector.DBConnAdapter = "logminer"
	config.Debezium.Connector.DBHistoryKafkaTopic = "history_topic_oracle"
	config.Debezium.Connector.Insensitive = "false"
	config.Debezium.Connector.LogMiningStrategy = "online_catalog"
	config.Debezium.Connector.SnapshotMode = "schema_only"

	tests := []struct {
		name    string
		wantErr bool
	}{
		{"test debezium connector setup", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := SetupConnector(); (err != nil) != tt.wantErr {
				t.Errorf("SetupConnector() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
