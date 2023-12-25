import { SupabaseConfig } from '@config/supabase';
import { Injectable } from '@nestjs/common';
import { SupabaseClient, createClient } from '@supabase/supabase-js';
import { SupabaseAuthClient } from '@supabase/supabase-js/dist/module/lib/SupabaseAuthClient';

// TODO: secure supabase connection
@Injectable()
export class SupabaseService {
  private readonly _supabase: SupabaseClient;

  constructor(config: SupabaseConfig) {
    this._supabase = createClient(config.url, config.serviceRoleKey);
  }

  auth(): SupabaseAuthClient {
    return this._supabase.auth;
  }
}
