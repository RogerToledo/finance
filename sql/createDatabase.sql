CREATE DATABASE Finance;
-----------------------------------------------------------------
CREATE TABLE credit_card (
	id uuid NOT NULL,
	"owner" varchar(50) NOT NULL,
	invoice_closing_day int4 NOT NULL,
	final_card_num int4 NULL,
	created_at timestamp DEFAULT now() NOT NULL,
	update_at timestamp NULL,
	CONSTRAINT credit_card_pk PRIMARY KEY (id),
	CONSTRAINT credit_card_unique UNIQUE (owner, final_card_num)
);
-----------------------------------------------------------------
CREATE TABLE payment_type (
	id uuid NOT NULL,
	"name" varchar(50) NOT NULL,
	create_at timestamp DEFAULT now() NOT NULL,
	update_at timestamp NULL,
	CONSTRAINT payment_type_pk PRIMARY KEY (id),
	CONSTRAINT payment_type_unique UNIQUE (name)
);
-----------------------------------------------------------------
CREATE TABLE person (
	id uuid NOT NULL,
	"name" varchar(50) NOT NULL,
	create_at timestamp DEFAULT now() NOT NULL,
	update_at timestamp NULL,
	CONSTRAINT person_pk PRIMARY KEY (id),
	CONSTRAINT person_unique UNIQUE (name)
);
-----------------------------------------------------------------
CREATE TABLE purchase_type (
	id uuid NOT NULL,
	"name" varchar(50) NOT NULL,
	create_at timestamp DEFAULT now() NOT NULL,
	update_at timestamp NULL,
	CONSTRAINT purchase_type_pk PRIMARY KEY (id),
	CONSTRAINT purchase_type_unique UNIQUE (name)
);
-----------------------------------------------------------------
CREATE TABLE purchase (
	id uuid NOT NULL,
	description varchar(150) NULL,
	amount numeric(10, 2) NOT NULL,
	"date" date NOT NULL,
	place varchar(60) NULL,
	paid bool NULL,
	id_payment_type uuid NOT NULL,
	id_purchase_type uuid NOT NULL,
	id_credit_card uuid NOT NULL,
	id_person uuid NOT NULL,
	create_at timestamp DEFAULT now() NOT NULL,
	update_at timestamp NULL,
	CONSTRAINT purchase_pk PRIMARY KEY (id)
);

-- finance.purchase foreign keys

ALTER TABLE purchase ADD CONSTRAINT purchase_credit_card_fk FOREIGN KEY (id_credit_card) REFERENCES credit_card(id) ON DELETE CASCADE ON UPDATE CASCADE;
ALTER TABLE purchase ADD CONSTRAINT purchase_payment_type_fk FOREIGN KEY (id_payment_type) REFERENCES payment_type(id) ON DELETE CASCADE ON UPDATE CASCADE;
ALTER TABLE purchase ADD CONSTRAINT purchase_person_fk FOREIGN KEY (id_person) REFERENCES person(id) ON DELETE CASCADE ON UPDATE CASCADE;
ALTER TABLE purchase ADD CONSTRAINT purchase_purchase_type_fk FOREIGN KEY (id_purchase_type) REFERENCES purchase_type(id) ON DELETE CASCADE ON UPDATE CASCADE;
-----------------------------------------------------------------
CREATE TABLE public.installment (
	id uuid NOT NULL,
	description varchar(50) NOT NULL,
	"month" date NOT NULL,
	value numeric(10, 2) NOT NULL,
	"number" int4 NOT NULL,
	paid bool NULL,
	purchase_id uuid NOT NULL,
	created_at timestamp DEFAULT now() not null,
	update_at timestamp NULL,
	CONSTRAINT installment_pk PRIMARY KEY (id)
);
-- installment foreign keys

ALTER TABLE installment ADD CONSTRAINT installment_purchase_fk FOREIGN KEY (purchase_id) REFERENCES purchase(id) ON DELETE CASCADE ON UPDATE CASCADE;
-----------------------------------------------------------------
