--
-- PostgreSQL database dump
--

-- Dumped from database version 16.3
-- Dumped by pg_dump version 16.3

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
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
-- Name: snippets; Type: TABLE; Schema: public; Owner: mnabil
--

CREATE TABLE public.snippets (
    id integer NOT NULL,
    title character varying(100) NOT NULL,
    content text NOT NULL,
    created_at timestamp without time zone NOT NULL,
    expires_at timestamp without time zone NOT NULL
);


ALTER TABLE public.snippets OWNER TO mnabil;

--
-- Name: snippets_id_seq; Type: SEQUENCE; Schema: public; Owner: mnabil
--

CREATE SEQUENCE public.snippets_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.snippets_id_seq OWNER TO mnabil;

--
-- Name: snippets_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: mnabil
--

ALTER SEQUENCE public.snippets_id_seq OWNED BY public.snippets.id;


--
-- Name: users; Type: TABLE; Schema: public; Owner: mnabil
--

CREATE TABLE public.users (
    id integer NOT NULL,
    name character varying(255) NOT NULL,
    email character varying(255) NOT NULL,
    password character(60) NOT NULL,
    created_at timestamp without time zone NOT NULL,
    active boolean DEFAULT true NOT NULL
);


ALTER TABLE public.users OWNER TO mnabil;

--
-- Name: users_id_seq; Type: SEQUENCE; Schema: public; Owner: mnabil
--

CREATE SEQUENCE public.users_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.users_id_seq OWNER TO mnabil;

--
-- Name: users_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: mnabil
--

ALTER SEQUENCE public.users_id_seq OWNED BY public.users.id;


--
-- Name: snippets id; Type: DEFAULT; Schema: public; Owner: mnabil
--

ALTER TABLE ONLY public.snippets ALTER COLUMN id SET DEFAULT nextval('public.snippets_id_seq'::regclass);


--
-- Name: users id; Type: DEFAULT; Schema: public; Owner: mnabil
--

ALTER TABLE ONLY public.users ALTER COLUMN id SET DEFAULT nextval('public.users_id_seq'::regclass);


--
-- Data for Name: snippets; Type: TABLE DATA; Schema: public; Owner: mnabil
--

COPY public.snippets (id, title, content, created_at, expires_at) FROM stdin;
1	An old silent pond	AN old silent pond...\\nA forg jumps into the pond,\\nsplash! Silence agan.\\nlash!	2024-06-10 17:49:21.657031	2025-06-10 17:49:21.657031
2	Second Entry	Second an old silent pond...\\nA forg jumps into the pond,\\nsplash! Silence agan.\\nlash!	2024-06-10 17:50:13.187779	2024-06-17 17:50:13.187779
4	Third entry	Hello\\nDarkness\\nMy old friend....	2024-06-19 09:56:32.899953	2024-06-26 09:56:32.899953
3	Third entry	Hello\\nDarkness\\nMy old friend....	2024-06-19 09:52:18.629967	2024-06-26 09:52:18.629967
5	My Personal Snippet	At first I was afraid, I was petrified\r\nThinking I could live without you by my side\r\nAnd after spending nights\r\nThinking how you did me wrong\r\nI grew strong\r\nAnd I learned how to get along	2024-06-20 12:00:58.988788	2024-06-21 12:00:58.988788
6	Gloria Gaylord - Survive	Go on, go, walk out the door\r\nTurn around now\r\nYou're not welcome anymore\r\nYou're the one who tried to hurt me with goodbye\r\nThink I'd crumble?\r\nYou think I'd lay down and die?	2024-06-20 12:25:40.028717	2025-06-20 12:25:40.028717
8	Lorem Ipsum	Lorem Ipsum is simply dummy text of the printing and typesetting industry. Lorem Ipsum has been the industry's standard dummy text ever since the 1500s, when an unknown printer took a galley of type and scrambled it to make a type specimen book. It has survived not only five centuries, but also the leap into electronic typesetting, remaining essentially unchanged. It was popularised in the 1960s with the release of Letraset sheets containing Lorem Ipsum passages, and more recently with desktop publishing software like Aldus PageMaker including versions of Lorem Ipsum.\r\n\r\nWhy do we use it?\r\nIt is a long established fact that a reader will be distracted by the readable content of a page when looking at its layout. The point of using Lorem Ipsum is that it has a more-or-less normal distribution of letters, as opposed to using 'Content here, content here', making it look like readable English. Many desktop publishing packages and web page editors now use Lorem Ipsum as their default model text, and a search for 'lorem ipsum' will uncover many web sites still in their infancy. Various versions have evolved over the years, sometimes by accident, sometimes on purpose (injected humour and the like).\r\n\r\n\r\nWhere does it come from?\r\nContrary to popular belief, Lorem Ipsum is not simply random text. It has roots in a piece of classical Latin literature from 45 BC, making it over 2000 years old. Richard McClintock, a Latin professor at Hampden-Sydney College in Virginia, looked up one of the more obscure Latin words, consectetur, from a Lorem Ipsum passage, and going through the cites of the word in classical literature, discovered the undoubtable source. Lorem Ipsum comes from sections 1.10.32 and 1.10.33 of "de Finibus Bonorum et Malorum" (The Extremes of Good and Evil) by Cicero, written in 45 BC. This book is a treatise on the theory of ethics, very popular during the Renaissance. The first line of Lorem Ipsum, "Lorem ipsum dolor sit amet..", comes from a line in section 1.10.32.\r\n\r\nThe standard chunk of Lorem Ipsum used since the 1500s is reproduced below for those interested. Sections 1.10.32 and 1.10.33 from "de Finibus Bonorum et Malorum" by Cicero are also reproduced in their exact original form, accompanied by English versions from the 1914 translation by H. Rackham.\r\n\r\nWhere can I get some?\r\nThere are many variations of passages of Lorem Ipsum available, but the majority have suffered alteration in some form, by injected humour, or randomised words which don't look even slightly believable. If you are going to use a passage of Lorem Ipsum, you need to be sure there isn't anything embarrassing hidden in the middle of text. All the Lorem Ipsum generators on the Internet tend to repeat predefined chunks as necessary, making this the first true generator on the Internet. It uses a dictionary of over 200 Latin words, combined with a handful of model sentence structures, to generate Lorem Ipsum which looks reasonable. The generated Lorem Ipsum is therefore always free from repetition, injected humour, or non-characteristic words etc.	2024-06-20 13:51:06.511604	2024-06-27 13:51:06.511604
9	New Snippet	This snippet is created at the time boredom ensues.	2024-06-20 17:02:46.735448	2024-06-27 17:02:46.735448
10	Example Lyrcis	test	2024-06-20 17:06:53.552117	2025-06-20 17:06:53.552117
11	My Homage To Peter Jackson	Peter Jackson,\r\n\r\nThe saviour of New Zealand\r\n\r\nThe King of Fantasy Filmmaking\r\n\r\nThe Greed of Fast Food...	2024-06-21 16:24:49.548039	2024-06-28 16:24:49.548039
\.


--
-- Data for Name: users; Type: TABLE DATA; Schema: public; Owner: mnabil
--

COPY public.users (id, name, email, password, created_at, active) FROM stdin;
1	Nabil	cucibaju123@gmail.com	$2a$12$LwvBfAFLgSfIasq3VKFcU.Z9w.vARImceLdhdDGD4PdyinUlisPVi	2024-06-21 12:59:48.636967	t
\.


--
-- Name: snippets_id_seq; Type: SEQUENCE SET; Schema: public; Owner: mnabil
--

SELECT pg_catalog.setval('public.snippets_id_seq', 11, true);


--
-- Name: users_id_seq; Type: SEQUENCE SET; Schema: public; Owner: mnabil
--

SELECT pg_catalog.setval('public.users_id_seq', 3, true);


--
-- Name: snippets snippets_pkey; Type: CONSTRAINT; Schema: public; Owner: mnabil
--

ALTER TABLE ONLY public.snippets
    ADD CONSTRAINT snippets_pkey PRIMARY KEY (id);


--
-- Name: users users_email_key; Type: CONSTRAINT; Schema: public; Owner: mnabil
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_email_key UNIQUE (email);


--
-- Name: users users_pkey; Type: CONSTRAINT; Schema: public; Owner: mnabil
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_pkey PRIMARY KEY (id);


--
-- PostgreSQL database dump complete
--

