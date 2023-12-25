import { IsNotEmpty, IsString } from 'class-validator';

export class SendPrivateMessageRequest {
  @IsString()
  @IsNotEmpty()
  receiverId!: string;

  @IsString()
  @IsNotEmpty()
  content!: string;
}
