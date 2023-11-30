import {
  CanActivate,
  ExecutionContext,
  Injectable,
  UnauthorizedException,
} from '@nestjs/common';
import { FastifyRequest } from 'fastify';
import { AuthService } from 'src/service/auth';

@Injectable()
export class VerifyConnectionGuard implements CanActivate {
  constructor(private readonly authService: AuthService) {}

  async canActivate(context: ExecutionContext): Promise<boolean> {
    const request = context.switchToHttp().getRequest<FastifyRequest>();
    const token = String(request.headers.authorization).replace('Bearer ', '');
    if (!token) {
      throw new UnauthorizedException('Token is required!');
    }

    // TODO: check if user connects to client service

    const user = await this.authService.getUser(token);
    request.user = user;
    return true;
  }
}
