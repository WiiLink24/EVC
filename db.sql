--
-- PostgreSQL database dump
--

\restrict Jov6I0aCzxXBahGXaBd1VSECWtPVMihQj3FBrr90FoRL1DeOK0eZSrzvpfeC9X9

-- Dumped from database version 18.3
-- Dumped by pg_dump version 18.3

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
    question_id integer NOT NULL,
    question_en text NOT NULL,
    question_de text NOT NULL,
    question_fr text NOT NULL,
    question_es text NOT NULL,
    question_it text NOT NULL,
    question_nl text NOT NULL,
    question_pt text NOT NULL,
    question_fr_ca text NOT NULL,
    response1_en text NOT NULL,
    response1_de text NOT NULL,
    response1_fr text NOT NULL,
    response1_es text NOT NULL,
    response1_it text NOT NULL,
    response1_nl text NOT NULL,
    response1_pt text NOT NULL,
    response1_fr_ca text NOT NULL,
    response2_en text NOT NULL,
    response2_de text NOT NULL,
    response2_fr text NOT NULL,
    response2_es text NOT NULL,
    response2_it text NOT NULL,
    response2_nl text NOT NULL,
    response2_pt text NOT NULL,
    response2_fr_ca text NOT NULL,
    extra bytea,
    category integer NOT NULL,
    date date NOT NULL,
    type character(1) NOT NULL,
    CONSTRAINT questions_type_check CHECK ((type = ANY (ARRAY['n'::bpchar, 'w'::bpchar])))
);


ALTER TABLE public.questions OWNER TO postgres;

--
-- Name: questions_question_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.questions_question_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.questions_question_id_seq OWNER TO postgres;

--
-- Name: questions_question_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.questions_question_id_seq OWNED BY public.questions.question_id;


--
-- Name: votes; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.votes (
    question_id integer NOT NULL,
    type_cd smallint NOT NULL,
    country_id integer NOT NULL,
    region_id integer NOT NULL,
    ans_cnt integer NOT NULL
);


ALTER TABLE public.votes OWNER TO postgres;

--
-- Name: questions question_id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.questions ALTER COLUMN question_id SET DEFAULT nextval('public.questions_question_id_seq'::regclass);


--
-- Data for Name: questions; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.questions (question_id, question_en, question_de, question_fr, question_es, question_it, question_nl, question_pt, question_fr_ca, response1_en, response1_de, response1_fr, response1_es, response1_it, response1_nl, response1_pt, response1_fr_ca, response2_en, response2_de, response2_fr, response2_es, response2_it, response2_nl, response2_pt, response2_fr_ca, extra, category, date, type) FROM stdin;
\.


--
-- Data for Name: votes; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.votes (question_id, type_cd, country_id, region_id, ans_cnt) FROM stdin;
\.


--
-- Name: questions_question_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.questions_question_id_seq', 1, false);


--
-- Name: questions questions_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.questions
    ADD CONSTRAINT questions_pkey PRIMARY KEY (question_id);


--
-- Name: votes votes_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.votes
    ADD CONSTRAINT votes_pkey PRIMARY KEY (question_id, type_cd, country_id, region_id);


--
-- Name: votes votes_question_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.votes
    ADD CONSTRAINT votes_question_id_fkey FOREIGN KEY (question_id) REFERENCES public.questions(question_id);


--
-- PostgreSQL database dump complete
--

\unrestrict Jov6I0aCzxXBahGXaBd1VSECWtPVMihQj3FBrr90FoRL1DeOK0eZSrzvpfeC9X9
