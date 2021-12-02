#### 避免通过时序攻击来枚举账户和密码

攻击者可以通过判断密码的比较时间，来重复的生成不同的密码进行大量的请求

```php
// bad 
class TimingSafeAuth
{
    private $db;
    
    public function __construct($db)
    {
		$this->db = $db;
        $this->dummy_pw = password_hash(random_bytes(32),PASSWORD_DEFAULT);
    }
    
    public function authenticate($username,$password)
    {
        $result = $this->db->get($user['username']);
        if ($result) {
            if (password_verify($password,$result['password'])) {
                return $result['id'];
            }
            return false;
        } else {
            // 下面这行return存在返回true的可能
			return password_verify($password,$this->dummy_pw); // x
            
            // 正确的写法
            // 为什么要多验证一下，是为了避免时序攻击，在去数据库取用户时，如果没有这个用户不进行，密码验证这个步骤的话，则返回时间与进行密码验证的时间是会有细微的不一致的，此时攻击者就可以根据时间的不一致进行攻击
            password_verify($password,$this->dummy_pw)
            return false
        }
    }
}
```

#### 加密和身份验证的区别

- **加密不保证数据的完整性**：一个被篡改了的数据仍然可以被解密，但是获取的数据就是错误的数据
- **身份验证不提供数据的机密性**：可以防止数据被篡改

#### 生产环境应当关闭

```ini
display_errors = Off
display_startup_errors = Off
error_reporting = E_ALL
log_errors = On
```

