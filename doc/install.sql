-- 短码库
CREATE TABLE sc_code on CLUSTER 2shards_2replicas
(

    `code` LowCardinality(String),

    `serial_number` Int64
    )
    ENGINE = Distributed('2shards_2replicas',
                         'sc',
                         'sc_rep_code',
                         javaHash(serial_number));


CREATE TABLE sc_rep_code on CLUSTER 2shards_2replicas
(

    `code` LowCardinality(String),

    `serial_number` Int64
    )
    ENGINE = ReplicatedMergeTree('/clickhouse/tables/{shard}/sc_v1/sc_rep_code',
                                 '{replica}')
    PARTITION BY toString(floor(serial_number / 20000000))
    PRIMARY KEY serial_number
    ORDER BY serial_number;



-- 当前已经使用的短码序号
CREATE TABLE sc_current_serial_number on CLUSTER 2shards_2replicas
(

    `serial_number` Int64,

    `create_time` Int64
)
    ENGINE = Distributed('2shards_2replicas',
                         'sc',
                         'sc_rep_current_serial_number',
                         javaHash(create_time));

CREATE TABLE sc_rep_current_serial_number on CLUSTER 2shards_2replicas
(

    `serial_number` Int64,

    `create_time` Int64
)
    ENGINE = ReplicatedMergeTree('/clickhouse/tables/{shard}/sc_v1/sc_rep_current_serial_number',
                                 '{replica}')
    PRIMARY KEY create_time
    ORDER BY create_time;

insert into sc_current_serial_number values(0,1718857660000);

-- 短码绑定的数据
CREATE TABLE sc_binding_data on CLUSTER 2shards_2replicas
(

    `code` LowCardinality(String),

    `message` String,

    `create_time` Int64
    )
    ENGINE = Distributed('2shards_2replicas',
                         'sc',
                         'sc_rep_binding_data',
                         javaHash(code));


CREATE TABLE sc_rep_binding_data on CLUSTER 2shards_2replicas
(

    `code` LowCardinality(String),

    `message` String,

    `create_time` Int64
    )
    ENGINE = ReplicatedMergeTree('/clickhouse/tables/{shard}/sc_v1/sc_rep_binding_data',
                                 '{replica}')
    PARTITION BY substring(code, 1,2)
    PRIMARY KEY code
    ORDER BY (code,
              create_time);
