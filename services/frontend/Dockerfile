# Use the official Node slim image as the base
FROM node:slim

# Set the working directory in the container
WORKDIR /app

# Copy package.json and package-lock.json first (if you have them)
# This allows for efficient layer caching when installing dependencies
COPY package*.json /app/

# Install dependencies (if you have them)
RUN npm install

RUN npm install -g nodemon

# Copy the rest of your application code into the container
COPY src /app

# Expose the port your app runs on (optional, helps for readability/documentation)
EXPOSE 3000

# Define the command to start your Node.js app
CMD ["nodemon", "index.js"]