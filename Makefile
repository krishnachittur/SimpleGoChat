NAME = app
DIR = chatapp
go:
	go build -o $(NAME) $(DIR)/*.go
clean:
	rm -f $(NAME)
