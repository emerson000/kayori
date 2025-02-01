# Use the official Node slim image as the base
FROM node:slim

# Set the working directory in the container
WORKDIR /app

# Copy package.json and package-lock.json first (if you have them)
# This allows for efficient layer caching when installing dependencies
COPY package*.json /app/

# Install dependencies (if you have them)
RUN npm install

# Copy the rest of your application code into the container
COPY . /app

# Ensure entrypoint.sh is copied into the container
COPY entrypoint.sh /app/entrypoint.sh

RUN npm run build

RUN chmod +x /app/entrypoint.sh

EXPOSE 3000

CMD /app/entrypoint.sh