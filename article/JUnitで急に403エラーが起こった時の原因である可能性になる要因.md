---
title: JUnitで急に403エラーが起こった時の原因である可能性になる要因
tags:
  - Java
  - spring-security
  - JUnit
  - SpringBoot
private: false
qiita_url: https://qiita.com/ShinguAkira/items/66deeeba8af76ae8e5ce
---

# JUnitで急に403エラーが起こった時の原因である可能性になる要因

恐らくSpring Boot 2.5以前のバージョンを使用していた場合に、2.7などのバージョンに変更した場合、Spring Securityによるテストが強制され、@AutoConfigureMockMvcで、
mockMvc.performを呼び出した際に、403となる。
直近でソースを修正していないにも関わらず、@SpringBootTestのテストが急にmockMvc.performを呼び出しているテストが403で想定結果と一致しない場合は、上記が原因の可能性が高い。
(私の場合は、Spring Boot 2.2から2.7への変更のため、当問題が発生しました。)

## 解決方法
- 引き続き、Spring Securityのテストを行う必要がない場合、

※Spring Securityを無視することになるので、バージョンアップによってSpring Securityのテストが追加で必要になる場合は使用しないでください。
(Spring Securityのロールを使用して、権限制御を行っている場合など、)
```diff_java
- @AutoConfigureMockMvc
+ @AutoConfigureMockMvc(addFilters = false)
```
(addFilters = false)を付与。
これにより、Spring SecurityのFilter Chainを無効化することができ、従来のテストを行うことができる。

- Spring Securityのテストを行う必要がある場合、
```diff_java
@Test
+ @WithMockUser(username = "admin", roles = {"ADMIN"})
void accessAdminEndpointWithAdminRole() throws Exception {
	mockMvc.perform(get("/admin-endpoint"))
		.andExpect(status().isOk());
}
```
@WithMockUseを使用して適切なロールを付与する。どのロールを設定するかは、ソース依存ですが、SecurityConfigなどの名前のファイルで、SecurityFilterChainを定義している部分で、
ロールによる制御を行っている可能性があるため、それに対応したロールを設定し、呼び出しを行うことでロールによる権限制御を含めたうえでテストを行うことができる。
※この場合、ADMINロールを設定してます。

特にSecurityFilterChainなど、SpringSecurityの設定を行っていない場合、デフォルトの設定が適用されます。
```diff_java
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;
import org.springframework.security.config.Customizer;
import org.springframework.security.config.annotation.web.configuration.EnableWebSecurity;
import org.springframework.security.web.SecurityFilterChain;

@Configuration
@EnableWebSecurity
public class SecurityConfig {

    @Bean
    public SecurityFilterChain filterChain(org.springframework.security.config.annotation.web.builders.HttpSecurity http) throws Exception {
        http
            // Authorize HTTP requests
            .authorizeHttpRequests(authorize -> authorize
                .anyRequest().authenticated() // All requests require authentication
            )
            // Enable form-based login with default settings
            .formLogin(Customizer.withDefaults())
            // Enable HTTP Basic authentication
            .httpBasic(Customizer.withDefaults())
            // Enable CSRF protection with default settings
            .csrf(Customizer.withDefaults())
            // Add default security headers
            .headers(Customizer.withDefaults());

        return http.build();
    }
}

```
ユーザーロールでのアクセス制御が設定されている場合は、以下のようにantMatchers("/user/**")hasRole("USER")が設定されています。

```diff_java
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;
import org.springframework.security.config.Customizer;
import org.springframework.security.config.annotation.web.builders.HttpSecurity;
import org.springframework.security.web.SecurityFilterChain;

@Configuration
public class SecurityConfig {
    
    @Bean
    public SecurityFilterChain filterChain(HttpSecurity http) throws Exception {
        http
            .authorizeHttpRequests(auth -> auth
                .antMatchers("/user/**").hasRole("USER")
                .anyRequest().authenticated()
            )
            .formLogin(Customizer.withDefaults())
            .httpBasic(Customizer.withDefaults());
        return http.build();
    }
}

```
