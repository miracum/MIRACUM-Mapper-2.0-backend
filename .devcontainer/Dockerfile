# Use the official Golang image as a base
FROM golang:1.22.2-alpine

# Install dependencies
RUN apk update && apk add --no-cache make git docker-compose zsh shadow

# Set the current working directory inside the container
WORKDIR /workspace

# Install Go tools
RUN go install golang.org/x/tools/gopls@latest
# Install Delve for debugging
RUN go install github.com/go-delve/delve/cmd/dlv@latest

# Set up a non-root user for VS Code
# RUN adduser -D vscode

# Change ownership of Go directories to the vscode user
# RUN chown -R vscode:vscode /go /usr/local/go/pkg

# Set zsh as the default shell for the vscode user
# RUN chsh -s /bin/zsh vscode
RUN chsh -s /bin/zsh root

# Copy the entrypoint script into the container
# COPY entrypoint.sh /usr/local/bin/entrypoint.sh
#
# Make the entrypoint script executable
# RUN chmod +x /usr/local/bin/entrypoint.sh

# Set the entrypoint to the script
# ENTRYPOINT ["/usr/local/bin/entrypoint.sh"]

# Switch to the vscode user
# USER vscode

# Expose port 8080
EXPOSE 8080