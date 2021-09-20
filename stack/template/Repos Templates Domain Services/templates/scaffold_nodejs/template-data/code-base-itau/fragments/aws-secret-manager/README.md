# AWS-SECRET-MANAGER

### Arquivos:

    app:
    - build.gradle
    - application.yml

### Perguntas:

1. Deseja configurar conexão com algum banco de dados? SIM


2. Qual o sistema de gerenciamento de banco de dados? (Opções)

    * MySql
      ***$CLI_AWS_DB_DRIVER*** = "com.amazonaws.secretsmanager.sql.AWSSecretsManagerMySQLDriver"
      ***$CLI_AWS_DB_DIALECT*** = "MYSQL"

    * PostgreSql
      ***$CLI_AWS_DB_DRIVER*** = "com.amazonaws.secretsmanager.sql.AWSSecretsManagerPostgreSQLDriver"
      ***$CLI_AWS_DB_DIALECT***    = "POSTGRE"

    * Oracle Database
      ***$CLI_AWS_DB_DRIVER*** = "com.amazonaws.secretsmanager.sql.AWSSecretsManagerOracleDriver"
      ***$CLI_AWS_DB_DIALECT***    = "ORACLE"


3. Informar a URL do banco de dados na AWS - ex. 'jdbc-secretsmanager:mysql:
   //database-mysql.crpx.us-east-1.rds.amazonaws.com:3306/testedb\' :
   $CLI_AWS_DB_URL


4. Informar a Secret Gerada na AWS para o banco de dados:
   $CLI_AWS_DB_USERNAME


5. Deseja configurar dados para resgatar o valor de uma Secret gerada no Secret Manager?

    * SIM
        - Informar a região da AWS: ***$CLI_AWS_SECRET_REGION*** (default = sa-east-1)

        - Informar a URL da região: ***$CLI_AWS_SECRET_URL*** (default = secretsmanager.sa-east-1.amazonaws.com)
