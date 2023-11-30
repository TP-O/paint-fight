import { Injectable, UnauthorizedException } from '@nestjs/common';
import { SupabaseService } from './supabase';
import { SupabaseAuthClient } from '@supabase/supabase-js/dist/module/lib/SupabaseAuthClient';

@Injectable()
export class AuthService {
  /**
   * Supabase authentication service.
   */
  private readonly _auth: SupabaseAuthClient;

  constructor(supabase: SupabaseService) {
    this._auth = supabase.auth();
  }

  /**
   * Get user by json web token.
   *
   * @param jwt The Supabase json web token.
   */
  async getUser(jwt: string) {
    const response = await this._auth.getUser(jwt).catch(() => {
      throw new UnauthorizedException('Invalid access token!');
    });

    if (response.error) {
      throw new UnauthorizedException(response.error.message);
    }

    return response.data.user;
  }
}
