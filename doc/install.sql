-- 短码库
CREATE TABLE short_code_code on CLUSTER short_code_2shards_2replicas
(

    `short_code` LowCardinality(String),

    `serial_number` Int64
    )
    ENGINE = Distributed('short_code_2shards_2replicas',
                         'short_code',
                         'short_code_rep_code',
                         javaHash(serial_number));

CREATE TABLE short_code_rep_code on CLUSTER short_code_2shards_2replicas
(

    `short_code` LowCardinality(String),

    `serial_number` Int64
    )
    ENGINE = ReplicatedMergeTree('/clickhouse/tables/{shard}/short_code/short_code_rep_code',
                                 '{replica}')
    PARTITION BY toString(floor(serial_number / 20000000))
    PRIMARY KEY serial_number
    ORDER BY serial_number;

-- 当前已经使用的短码序号
CREATE TABLE short_code_current_serial_number on CLUSTER short_code_2shards_2replicas
(

    `serial_number` Int64,

    `create_time` Int64
)
    ENGINE = Distributed('short_code_2shards_2replicas',
                         'short_code',
                         'short_code_rep_current_serial_number',
                         javaHash(create_time));

CREATE TABLE short_code_rep_current_serial_number on CLUSTER short_code_2shards_2replicas
(

    `serial_number` Int64,

    `create_time` Int64
)
    ENGINE = ReplicatedMergeTree('/clickhouse/tables/{shard}/short_code/short_code_rep_current_serial_number',
                                 '{replica}')
    PRIMARY KEY create_time
    ORDER BY create_time;

-- 短码绑定的数据
CREATE TABLE short_code_binding_data on CLUSTER short_code_2shards_2replicas
(

    `short_code` LowCardinality(String),

    `message` String,

    `create_time` Int64
    )
    ENGINE = Distributed('short_code_2shards_2replicas',
                         'short_code',
                         'short_code_rep_binding_data',
                         javaHash(short_code));


CREATE TABLE short_code_rep_binding_data on CLUSTER short_code_2shards_2replicas
(

    `short_code` LowCardinality(String),

    `message` String,

    `create_time` Int64
    )
    ENGINE = ReplicatedMergeTree('/clickhouse/tables/{shard}/short_code/short_code_rep_binding_data',
                                 '{replica}')
    PARTITION BY substring(short_code, 1,2)
    PRIMARY KEY short_code
    ORDER BY (short_code,
              create_time);