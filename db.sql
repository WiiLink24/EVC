--
-- PostgreSQL database dump
--

\restrict 1QmxivAeLIsEPdjMpxfAgTdldhC86Pkj233Sdy3LqlQu2J9Kh1RQCbvbfmTz1Um

-- Dumped from database version 18.4
-- Dumped by pg_dump version 18.4

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET transaction_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET xmloption = content;
SET client_min_messages = warning;
SET row_security = off;

SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- Name: questions; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.questions (
    question_id bigint,
    content_english character varying(255),
    content_german character varying(255),
    content_french character varying(255),
    content_spanish character varying(255),
    content_italian character varying(255),
    content_dutch character varying(255),
    content_portuguese character varying(255),
    content_french_canada character varying(255),
    choice1_english character varying(255),
    choice1_german character varying(255),
    choice1_french character varying(255),
    choice1_spanish character varying(255),
    choice1_italian character varying(255),
    choice1_dutch character varying(255),
    choice1_portuguese character varying(255),
    choice1_french_canada character varying(255),
    choice2_english character varying(255),
    choice2_german character varying(255),
    choice2_french character varying(255),
    choice2_spanish character varying(255),
    choice2_italian character varying(255),
    choice2_dutch character varying(255),
    choice2_portuguese character varying(255),
    choice2_french_canada character varying(255),
    type character varying(1),
    category smallint,
    date date,
    end_date date
);


ALTER TABLE public.questions OWNER TO postgres;

--
-- Name: suggestions; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.suggestions (
    id integer NOT NULL,
    country_code integer,
    region_code integer,
    language_code integer,
    content character varying,
    choice1 character varying,
    choice2 character varying,
    wii_no bigint
);


ALTER TABLE public.suggestions OWNER TO postgres;

--
-- Name: suggestions_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.suggestions_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.suggestions_id_seq OWNER TO postgres;

--
-- Name: suggestions_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.suggestions_id_seq OWNED BY public.suggestions.id;


--
-- Name: user; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public."user" (
    id integer NOT NULL,
    username character varying(100),
    password_hash character varying
);


ALTER TABLE public."user" OWNER TO postgres;

--
-- Name: user_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.user_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.user_id_seq OWNER TO postgres;

--
-- Name: user_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.user_id_seq OWNED BY public."user".id;


--
-- Name: votes; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.votes (
    id bigint,
    type_cd smallint,
    question_id bigint,
    wii_no bigint,
    country_id smallint,
    region_id smallint,
    ans_cnt bigint,
    new_column integer NOT NULL
);


ALTER TABLE public.votes OWNER TO postgres;

--
-- Name: votes_new_column_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.votes_new_column_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.votes_new_column_seq OWNER TO postgres;

--
-- Name: votes_new_column_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.votes_new_column_seq OWNED BY public.votes.new_column;


--
-- Name: suggestions id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.suggestions ALTER COLUMN id SET DEFAULT nextval('public.suggestions_id_seq'::regclass);


--
-- Name: user id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public."user" ALTER COLUMN id SET DEFAULT nextval('public.user_id_seq'::regclass);


--
-- Name: votes id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.votes ALTER COLUMN id SET DEFAULT nextval('public.votes_new_column_seq'::regclass);


--
-- Name: votes new_column; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.votes ALTER COLUMN new_column SET DEFAULT nextval('public.votes_new_column_seq'::regclass);


--
-- Name: user user_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public."user"
    ADD CONSTRAINT user_pkey PRIMARY KEY (id);


--
-- PostgreSQL database dump complete
--

\unrestrict 1QmxivAeLIsEPdjMpxfAgTdldhC86Pkj233Sdy3LqlQu2J9Kh1RQCbvbfmTz1Um

