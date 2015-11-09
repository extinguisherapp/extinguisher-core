require_relative './bundle/bundler/setup'
require "socket"
require "compass"
require "compass/exec"


class SocketServer

  attr_accessor :path

  def initialize(path:, socket_path:)

    @path = path
    File.unlink(socket_path) if File.exists?(socket_path)
    @socket_path = socket_path

    @server ||= UNIXServer.new(@socket_path)
  end

  def start
    while true
      socket = @server.accept
      Compass::Commands::UpdateProject.new('./', {}).perform
      socket.close
    end
  end
end

Dir.chdir(ARGV[0])
Compass::Commands::UpdateProject.new(ARGV[0], {}).perform
SocketServer.new(path: ARGV[0], socket_path: File.join(ARGV[0],"css.socket")).start
