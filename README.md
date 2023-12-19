# camel

Camel is a command-line interface (CLI) used to locally interact with llama2, an AI from Meta.

## Get Started

Check [get started video on youtube here](https://youtu.be/cs-aEjyixGU?si=X5sBF2wz_RbKeUiZ).

1. Run the Ollama Docker container:

```shell
sudo docker run -d -v ollama:/root/.ollama -p 11434:11434 --name ollama ollama/ollama
```

For more detailed information, refer to the [Ollama Quickstart Docker](https://hub.docker.com/r/ollama/ollama). Please note we are using CPU only, the AI will response slow, if you have GPU, you can follow the instruction to run the docker and using your GPU to improve performance.

2. Pull the llama2 model:

```shell
curl --location 'http://localhost:11434/api/pull' \
--header 'Content-Type: application/json' \
--data '{
    "name": "llama2:7b"
}'
```

3. Download camel binary from [release page](https://github.com/Hidayathamir/camel/releases).

4. Run camel binary:

```shell
./camel
```

Upon the initial run, you'll receive an information message as follows. Camel will create a folder named `camel_data` and a file named `camel_data/question.md`. The message will appear like this:

```shell
{"level":"info","msg":"please write your question in `camel_data/question.md` file","time":"2023-12-19T22:34:09+07:00"}
```

Write your question in the `camel_data/question.md` file.

After writing your question, execute the camel binary again. It will read your question from `camel_data/question.md`.

The response will be displayed in the console and saved in the file `camel_data/answer.md`. Your chat history will also be saved in `camel_data/history.json`.

Simply add your question to `camel_data/question.md`, and run the Camel binary whenever you wish to get a response.
