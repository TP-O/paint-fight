import { Injectable } from '@nestjs/common';
import { SupabaseClient, createClient } from '@supabase/supabase-js';
import { SupabaseAuthClient } from '@supabase/supabase-js/dist/module/lib/SupabaseAuthClient';
import { SupabaseConfig } from 'src/config/supabase';

// TODO: secure supabase connection
@Injectable()
export class SupabaseService {
  private supabase: SupabaseClient;

  constructor(config: SupabaseConfig) {
    this.supabase = createClient(config.url, config.serviceRoleKey);
  }

  auth(): SupabaseAuthClient {
    return this.supabase.auth;
  }
}
