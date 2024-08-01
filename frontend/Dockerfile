# Stage 1: Build the React app
FROM node:16-alpine AS build

# Set working directory
WORKDIR /app

# Install dependencies
COPY package.json package-lock.json ./
RUN npm install

# Copy source files and build the app
COPY . ./
RUN npm run build

# Stage 2: Serve the app using nginx
FROM nginx:stable-alpine

# Copy the built React app from the previous stage to nginx's html directory
COPY --from=build /app/dist /usr/share/nginx/html

# Copy custom nginx configuration file
COPY nginx.conf /etc/nginx/conf.d/default.conf


# Start nginx server
CMD ["nginx", "-g", "daemon off;"]