CREATE TABLE transactions (
	uuid uuid PRIMARY KEY,
	created_at timestamp NOT NULL,
	updated_at timestamp NOT NULL,
	request_local_id varchar(255),
	request_service varchar(255),
	request_phone varchar(255),
	request_amount varchar(255),
	status varchar(255),
	error_code int,
	error_msg varchar(255),
	result_status varchar(255),
	result_ref_num bigint,
	result_service varchar(255),
	result_destination varchar(255),
	result_amount bigint,
	result_state varchar(255),
	is_checked boolean NOT NULL DEFAULT false,
	client varchar(255)
);

CREATE TABLE users (
	username varchar(255) PRIMARY KEY,
	created_at timestamp NOT NULL,
	updated_at timestamp NOT NULL,
	password varchar(255) NOT NULL
);